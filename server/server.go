package server

import (
	"context"
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/config"
	"github.com/amanbolat/ca-warehouse-client/filemaker"
	"github.com/amanbolat/ca-warehouse-client/printing"
	"github.com/amanbolat/gofmcon"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Server struct {
	router        *echo.Echo
	memCache      *cache.Cache
	logger        *logrus.Logger
	wsServer      *melody.Melody
	wsSessions    *sync.Map
	shipmentStore *filemaker.ShipmentStore
}

func NewServer(config config.Config, logger *logrus.Logger) Server {
	conn := gofmcon.NewFMConnector(config.FmHost, "", config.FmUser, config.FmPass)
	entryStore := filemaker.NewEntryStore(conn, config.FmDatabaseName)
	entryStore.FMConn().SetDebug(config.Debug)
	shipmentStore := filemaker.NewShipmentStore(conn, config.FmDatabaseName)
	customerStore := filemaker.NewCustomerStore(conn, config.FmDatabaseName)

	s := Server{
		memCache:      cache.New(time.Hour*24, time.Hour*30),
		logger:        logger,
		wsServer:      melody.New(),
		wsSessions:    &sync.Map{},
		shipmentStore: shipmentStore,
	}

	e := echo.New()
	if config.Debug {
		e.Debug = true
	}

	lm, err := printing.NewLabelManger(config.FontPath)
	if err != nil {
		log.Fatal(err)
	}

	var a = API{
		entryStore:       entryStore,
		shipmentStore:    shipmentStore,
		customerStore:    customerStore,
		memCache:         cache.New(time.Minute*5, time.Minute*7),
		kdniaoApi:        api.NewKDNiaoApi(config.KDNiaoConfig),
		printer:          printing.Printer{Name: config.Printer},
		labelManager:     lm,
		apiRequestsCache: s.memCache,
	}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.Logger().Error(err)
		apiErr, ok := err.(api.Error)
		if ok {
			err = c.JSON(http.StatusServiceUnavailable, apiErr)
			if err != nil {
				c.Logger().Error(err)
			}
		} else {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}",` +
			`"method":"${method}","uri":"${uri}","status":${status},` +
			`"latency_human":"${latency_human}"` + "\n",
		Output: os.Stdout,
	}), middleware.Recover(), middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:80", "https://wh.me", "http://wh.me"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, "X-API-REQUEST-ID"},
		AllowCredentials: true,
	}))

	g := e.Group("/api")
	g.GET("/entries", a.GetEntryList)
	g.GET("/entries/:id", a.GetEntrySingle)
	g.POST("/entries/:id/print_barcode", a.PrintEntryBarcode)
	g.POST("/entries", s.duplicatePreventMiddleware(a.CreateEntry))
	g.PATCH("/entries", a.EditEntry)
	g.GET("/shipments", a.GetShipmentList)
	g.GET("/shipments/:code", a.GetShipmentSingle)
	g.POST("/shipments/:code/print/unit_loads", a.PrintShipmentULLabels)
	g.POST("/shipments/:code/print/preparation_info", a.PrintShipmentPreparationInfo)
	g.POST("/shipments/:code/print/partner_info", a.PrintShipmentPartnerInfo)
	g.GET("/customers", a.GetCustomerList)
	g.GET("/kdniao/get_source/:track_code", a.GetSourceByTrackCode)

	s.router = e
	s.SetupWsServer()
	s.StartShipmentUpdates()

	return s
}

func (s Server) Start(port int) {
	go func() {
		if err := s.router.Start(fmt.Sprintf(":%d", port)); err != nil {
			s.router.Logger.Fatalf("failed to start server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.router.Shutdown(ctx); err != nil {
		s.router.Logger.Fatal(err)
	}
}

func (s Server) SetupWsServer() {
	s.wsServer.HandleConnect(func(sess *melody.Session) {
		s.wsSessions.Load(sess)
		s.logger.Info("New ws user connection")
	})
	s.wsServer.HandleDisconnect(func(sess *melody.Session) {
		s.wsSessions.Delete(sess)
		s.logger.Info("Ws user disconnected")
	})
	wsGroup := s.router.Group("/ws")
	wsGroup.GET("/shipments", func(c echo.Context) error {
		return s.wsServer.HandleRequest(c.Response(), c.Request())
	})
}

func (s Server) duplicatePreventMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Method != http.MethodPost {
			c.Logger().Warnf("duplicatePreventMiddleware used on method %s, buy should be used on POST only", c.Request().Method)
			return next(c)
		}
		apiRequestId := c.Request().Header.Get(XApiRequestId)
		if apiRequestId == "" {
			return api.NewError(echo.ErrServiceUnavailable, "没有收到请求ID", "")
		}

		_, ok := s.memCache.Get(apiRequestId)
		if ok {
			return api.NewError(echo.ErrServiceUnavailable, "不允许重复请求", "建议您刷新页面再试试")
		}

		s.memCache.SetDefault(apiRequestId, true)
		c.Set(XApiRequestId, apiRequestId)

		return next(c)
	}
}
