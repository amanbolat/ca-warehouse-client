package server

import (
	"context"
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/config"
	"github.com/amanbolat/ca-warehouse-client/filemaker"
	model "github.com/amanbolat/ca-warehouse-client/logistics"
	"github.com/amanbolat/ca-warehouse-client/printing"
	"github.com/amanbolat/gofmcon"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/patrickmn/go-cache"
	"gopkg.in/olahol/melody.v1"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Server struct {
	router *echo.Echo
}

func NewServer(config config.Config) Server {
	e := echo.New()
	if config.Debug {
		e.Debug = true
	}
	conn := gofmcon.NewFMConnector(config.FmHost, "", config.FmUser, config.FmPass)
	entryStore := filemaker.NewEntryStore(conn, config.FmDatabaseName)
	shipmentStore := filemaker.NewShipmentStore(conn, config.FmDatabaseName)
	customerStore := filemaker.NewCustomerStore(conn, config.FmDatabaseName)

	lm, err := printing.NewLabelManger(config.FontPath)
	if err != nil {
		log.Fatal(err)
	}

	m := NewBroadcaster()
	var a = API{
		entryStore:    entryStore,
		shipmentStore: shipmentStore,
		customerStore: customerStore,
		wsServer:      m,
		memCache:      cache.New(time.Minute*5, time.Minute*7),
		kdniaoApi:     api.NewKDNiaoApi(config.KDNiaoConfig),
		printer:       printing.Printer{Name: config.Printer},
		labelManager:  lm,
	}
	e.Use(middleware.Logger(), middleware.Recover(), middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:80", "https://wh.me", "http://wh.me"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	g := e.Group("/api")
	g.GET("/entries", a.GetEntryList)
	g.GET("/entries/:id", a.GetEntrySingle)
	g.POST("/entries/:id/print_barcode", a.PrintEntryBarcode)
	g.POST("/entries", a.CreateEntry)
	g.PATCH("/entries", a.EditEntry)
	g.GET("/shipments", a.GetShipmentList)
	g.GET("/shipments/:code", a.GetShipmentSingle)
	g.POST("/shipments/:code/print/unit_loads", a.PrintShipmentULLabels)
	g.POST("/shipments/:code/print/preparation_info", a.PrintShipmentPreparationInfo)
	g.POST("/shipments/:code/print/partner_info", a.PrintShipmentPartnerInfo)
	g.GET("/customers", a.GetCustomerList)
	g.GET("/kdniao/get_source/:track_code", a.GetSourceByTrackCode)

	s := Server{
		router: e,
	}

	return s
}

func (s Server) Start(port int) {
	go func() {
		if err := s.router.Start(fmt.Sprintf(":%d", port)); err != nil {
			s.router.Logger.Info("shutting down the server")
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

type BroadcastingChans struct {
	Shipments []model.Shipment
}

func NewBroadcaster() *melody.Melody {
	m := melody.New()

	// m.HandleConnect()
	//
	// m.HandleDisconnect()

	return m
}
