package server

import (
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/filemaker"
	"github.com/amanbolat/ca-warehouse-client/printing"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"os/exec"
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

	meta = warehouse.MapEntryFields(meta)

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

func (a API) PrintEntryBarcode(c echo.Context) error {
	entryId := c.Param("id")

	bm := printing.BarcodeManger{}

	bc, err := bm.CreateEntryBarcode(entryId)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	printCmd := exec.Command("lpr", "-P", "Canon_G3000_series", "-o", "media=a4", "-r", bc.FullPath)
	out, err := printCmd.CombinedOutput()
	if err != nil {
		c.Logger().Errorf("%v: %s", err, string(out))
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, "done")
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
