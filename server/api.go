package server

import (
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/filemaker"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"gopkg.in/olahol/melody.v1"
	"net/http"
)

type JSONResponse struct {
	Meta api.ResponseMeta `json:"meta"`
	Data interface{}      `json:"data"`
}

type API struct {
	entryStore    *filemaker.EntryStore
	shipmentStore *filemaker.ShipmentStore
	wsServer      *melody.Melody
}

func (a API) GetEntryList(c echo.Context) error {
	meta := api.RequestMeta{}
	d := schema.NewDecoder()
	err := d.Decode(&meta, c.QueryParams())
	if err != nil {
		return echo.ErrBadRequest
	}

	entries, res, err := a.entryStore.GetEntryList(meta)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, JSONResponse{
		Meta: res,
		Data: entries,
	})
}

func (a API) GetEntrySingle(c echo.Context) error {
	return nil
}

func (a API) GetShipmentList(c echo.Context) error {
	meta := api.RequestMeta{}
	d := schema.NewDecoder()
	err := d.Decode(&meta, c.QueryParams())
	if err != nil {
		return echo.ErrBadRequest
	}

	shipments, res, err := a.shipmentStore.GetShipmentList(meta)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, JSONResponse{
		Meta: res,
		Data: shipments,
	})
}

func (a API) EditEntry(c echo.Context) error {
	return nil
}

func (a API) CreateEntry(c echo.Context) error {
	return nil
}

func (a API) ShipmentsWS(c echo.Context) error {
	return a.wsServer.HandleRequest(c.Response(), c.Request())
}
