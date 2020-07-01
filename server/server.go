package server

import (
	"context"
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/config"
	"github.com/amanbolat/ca-warehouse-client/filemaker"
	model "github.com/amanbolat/ca-warehouse-client/logistics"
	"github.com/amanbolat/gofmcon"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/olahol/melody.v1"
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

	m := NewBroadcaster()
	a := API{entryStore: entryStore, shipmentStore: shipmentStore, wsServer: m}

	g := e.Group("/api", middleware.Logger(), middleware.Recover())
	g.GET("/entries", a.GetEntryList)
	g.GET("/entries/{}", a.GetEntrySingle)
	g.POST("/entries", a.CreateEntry)
	g.PUT("/entries", a.EditEntry)
	g.GET("/shipments", a.GetShipmentList)

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
