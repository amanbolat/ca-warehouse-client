package server

import (
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/crm"
	"github.com/amanbolat/ca-warehouse-client/filemaker"
	"github.com/amanbolat/ca-warehouse-client/printing"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"strconv"
)

type JSONResponse struct {
	Meta api.ResponseMeta `json:"meta"`
	Data interface{}      `json:"data"`
}

type API struct {
	entryStore    *filemaker.EntryStore
	shipmentStore *filemaker.ShipmentStore
	customerStore *filemaker.CustomerStore
	wsServer      *melody.Melody
	memCache      *cache.Cache
	kdniaoApi     *api.KDNiaoApi
	printer       printing.Printer
	labelManager  printing.LabelManager
}

func (a API) GetEntryList(c echo.Context) error {
	meta := api.RequestMeta{}
	d := schema.NewDecoder()
	err := d.Decode(&meta, c.QueryParams())
	if err != nil {
		return echo.ErrBadRequest
	}

	meta = warehouse.MapEntryFields(meta)
	meta.InternalFilter["Warehouse"] = "=GZWH2"
	meta.InternalFilter["Id_shipmentNumber"] = "="
	meta.InternalFilter["is_utilized"] = "="

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

func (a API) GetShipmentSingle(c echo.Context) error {
	code := c.Param("code")

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, JSONResponse{
		Meta: api.ResponseMeta{Page: 1, Total: 1, Count: 1},
		Data: sm,
	})
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

	bc, err := a.labelManager.CreateEntryBarcode(entryId)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	err = a.printer.PrintFiles(1, "", bc.FullPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, "done")
}

func (a API) GetCustomerList(c echo.Context) error {
	inMemCustomers, ok := a.memCache.Get("customer_list")
	if ok {
		c.Logger().Print("Got customer from cache")
		customers := inMemCustomers.([]crm.Customer)
		total := len(customers)
		return c.JSON(http.StatusOK, JSONResponse{
			Meta: api.ResponseMeta{
				Count: total,
				Total: total,
			},
			Data: customers,
		})
	}

	meta := api.RequestMeta{}
	d := schema.NewDecoder()
	err := d.Decode(&meta, c.QueryParams())
	if err != nil {
		return echo.ErrBadRequest
	}

	customers, res, err := a.customerStore.GetCustomerList(meta)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	a.memCache.SetDefault("customer_list", customers)

	return c.JSON(http.StatusOK, JSONResponse{
		Meta: res,
		Data: customers,
	})
}

func (a API) GetSourceByTrackCode(c echo.Context) error {
	trackCode := c.Param("track_code")

	kdniaoResponse, err := a.kdniaoApi.GetSourceByTrack(trackCode)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	return c.JSON(http.StatusOK, kdniaoResponse)
}

func (a API) PrintShipmentULLabels(c echo.Context) error {
	code := c.Param("code")

	copies, _ := strconv.Atoi(c.QueryParam("copies"))

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	label, err := a.labelManager.CreateUnitLoadLabels(sm)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	err = a.printer.PrintFiles(copies, "", label.FullPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, "done")
}

func (a API) PrintShipmentPreparationInfo(c echo.Context) error {
	code := c.Param("code")

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	l, err := a.labelManager.CreateShipmentPreparationLabels(sm)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	err = a.printer.PrintFiles(1, "", l.FullPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, "done")
}

func (a API) PrintShipmentPartnerInfo(c echo.Context) error {
	code := c.Param("code")

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrNotFound
	}

	l, err := a.labelManager.CreateShipmentPartnerInfoLabel(sm)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}
	c.Logger().Printf("printing shipmentPartnerInfo, path to file: %s", l.FullPath)

	err = a.printer.PrintFiles(1, "", l.FullPath)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return c.String(http.StatusOK, "done")
}

func (a API) EditEntry(c echo.Context) error {
	return nil
}

func (a API) CreateEntry(c echo.Context) error {
	entry := &warehouse.Entry{}
	err := c.Bind(entry)

	if err != nil {
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	fmt.Println(entry)

	newEntry, err := a.entryStore.CreateEntry(*entry)
	if err != nil {
		c.Logger().Error(err)
		return echo.ErrBadRequest
	}

	return c.JSON(http.StatusOK, newEntry)
}

func (a API) ShipmentsWS(c echo.Context) error {
	return a.wsServer.HandleRequest(c.Response(), c.Request())
}
