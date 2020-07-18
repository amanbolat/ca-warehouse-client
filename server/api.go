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
	"net/http"
	"strconv"
)

type JSONResponse struct {
	Meta api.ResponseMeta `json:"meta"`
	Data interface{}      `json:"data"`
}

type API struct {
	entryStore       *filemaker.EntryStore
	shipmentStore    *filemaker.ShipmentStore
	customerStore    *filemaker.CustomerStore
	memCache         *cache.Cache
	kdniaoApi        *api.KDNiaoApi
	printer          printing.Printer
	labelManager     printing.LabelManager
	apiRequestsCache *cache.Cache
}

// XApiRequestId used to prevent duplicated POST requests
const XApiRequestId = "X-API-REQUEST-ID"

var singleRecordMeta = api.ResponseMeta{
	Page:  1,
	Count: 1,
	Total: 1,
}

func (a API) GetEntryList(c echo.Context) error {
	meta := api.RequestMeta{}
	d := schema.NewDecoder()
	err := d.Decode(&meta, c.QueryParams())
	if err != nil {
		return api.NewError(err, "请求有误", "建议您联系管理员")
	}

	meta = warehouse.MapEntryFields(meta)
	meta.InternalFilter["Warehouse"] = "=GZWH2"
	meta.InternalFilter["Id_shipmentNumber"] = "="
	meta.InternalFilter["is_utilized"] = "="

	entries, res, err := a.entryStore.GetEntryList(meta)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, JSONResponse{
		Meta: res,
		Data: entries,
	})
}

func (a API) GetEntrySingle(c echo.Context) error {
	id := c.Param("id")
	e, err := a.entryStore.GetEntryById(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, JSONResponse{
		Meta: singleRecordMeta,
		Data: e,
	})
}

func (a API) GetShipmentSingle(c echo.Context) error {
	code := c.Param("code")

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, JSONResponse{
		Meta: singleRecordMeta,
		Data: sm,
	})
}

func (a API) GetShipmentList(c echo.Context) error {
	meta := api.RequestMeta{}
	d := schema.NewDecoder()
	err := d.Decode(&meta, c.QueryParams())
	if err != nil {
		return api.NewError(err, "请求有误", "建议您联系管理员")
	}

	shipments, res, err := a.shipmentStore.GetShipmentList(meta)
	if err != nil {
		return err
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
		return api.NewError(err, "无法生成入库标签", "建议您联系管理员")
	}

	err = a.printer.PrintFiles(1, "", bc.FullPath)
	if err != nil {
		return api.NewError(err, "无法打印入库标签", "建议您联系管理员")
	}

	return c.String(http.StatusOK, "done")
}

func (a API) GetCustomerList(c echo.Context) error {
	inMemCustomers, ok := a.memCache.Get("customer_list")
	if ok {
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
		return api.NewError(err, "请求有误", "建议您联系管理员")
	}

	customers, res, err := a.customerStore.GetCustomerList(meta)
	if err != nil {
		return err
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
		return api.NewError(err, fmt.Sprintf("没有找到 %s 快递单的来源", trackCode), "")
	}

	return c.JSON(http.StatusOK, kdniaoResponse)
}

func (a API) PrintShipmentULLabels(c echo.Context) error {
	code := c.Param("code")

	copies, _ := strconv.Atoi(c.QueryParam("copies"))

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		return err
	}

	label, err := a.labelManager.CreateUnitLoadLabels(sm)
	if err != nil {
		return api.NewError(err, "生成货物标签遇到错误", "建议您联系管理员")
	}

	err = a.printer.PrintFiles(copies, "", label.FullPath)
	if err != nil {
		return api.NewError(err, "打印货物标签遇到错误", "建议您联系管理员")
	}

	return c.String(http.StatusOK, "done")
}

func (a API) PrintShipmentPreparationInfo(c echo.Context) error {
	code := c.Param("code")

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		return err
	}

	l, err := a.labelManager.CreateShipmentPreparationLabels(sm)
	if err != nil {
		return api.NewError(err, "无法生成发货明细", "建议您联系管理员")
	}

	err = a.printer.PrintFiles(1, "", l.FullPath)
	if err != nil {
		return api.NewError(err, "打印货物明细遇到错误", "建议您联系管理员")
	}

	return c.String(http.StatusOK, "done")
}

func (a API) PrintShipmentPartnerInfo(c echo.Context) error {
	code := c.Param("code")

	sm, err := a.shipmentStore.GetShipmentByCode(code)
	if err != nil {
		return err
	}

	l, err := a.labelManager.CreateShipmentPartnerInfoLabel(sm)
	if err != nil {
		return api.NewError(err, "无法生成合作方货物明细", "建议您联系管理员")
	}

	err = a.printer.PrintFiles(1, "", l.FullPath)
	if err != nil {
		return api.NewError(err, "打印合作方货物明细遇到错误", "建议您联系管理员")
	}

	return c.String(http.StatusOK, "done")
}

func (a API) EditEntry(c echo.Context) error {
	e := &warehouse.Entry{}
	err := c.Bind(e)
	if err != nil {
		return api.NewError(err, "请求有误", "请核对信息或者联系管理员")
	}

	updatedEntry, err := a.entryStore.UpdateEntry(*e)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, updatedEntry)
}

func (a API) CreateEntry(c echo.Context) error {
	entry := &warehouse.Entry{}
	err := c.Bind(entry)

	if err != nil {
		a.removeApiRequestId(c)
		return api.NewError(err, "请求有误", "有可能新加的入库数据有误。建议您联系管理员")
	}

	newEntry, err := a.entryStore.CreateEntry(*entry)
	if err != nil {
		a.removeApiRequestId(c)
		return err
	}

	return c.JSON(http.StatusOK, newEntry)
}

// removeApiRequestId removes XApiRequestId from
// apiRequestsCache. Used when request is failed
func (a API) removeApiRequestId(c echo.Context) {
	apiRequestId, ok := c.Get(XApiRequestId).(string)
	if ok {
		a.apiRequestsCache.Delete(apiRequestId)
	}
}
