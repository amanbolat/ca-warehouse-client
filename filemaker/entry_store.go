package filemaker

import (
	"encoding/json"
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/filemaker/fmutil"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	fm "github.com/amanbolat/gofmcon"
	"github.com/pkg/errors"
	"strconv"
)

var ErrZeroRecordsInResultSet = errors.New("0 records in fmResultSet after entry update")

const (
	ENTRY_LAYOUT = "warehouse_entry_single"
)

type EntryStore struct {
	conn         *fm.FMConnector
	databaseName string
}

func (s *EntryStore) DBName() string {
	return s.databaseName
}

func (s *EntryStore) FMConn() *fm.FMConnector {
	return s.conn
}

func NewEntryStore(conn *fm.FMConnector, dbName string) *EntryStore {
	return &EntryStore{
		conn:         conn,
		databaseName: dbName,
	}
}

func (s *EntryStore) GetEntryById(id string) (warehouse.Entry, error) {
	q := fm.NewFMQuery(s.databaseName, ENTRY_LAYOUT, fm.Find)
	q.WithFields(
		fm.FMQueryField{Name: "id", Value: id, Op: fm.Equal},
	)

	rec, err := fmutil.GetFileMakerRecordSingle(s, q)

	if err != nil {
		return warehouse.Entry{}, api.NewError(err, fmt.Sprintf("没有找到id为 %s 的入库", id), "原因无知，请联系管理员")
	}

	fEntry := warehouse.FileMakerEntry{}
	b, err := rec.JsonFields()
	if err != nil {
		return warehouse.Entry{}, api.NewError(err, fmt.Sprintf("没有找到id为 %s 的入库", id), "原因无知，请联系管理员")
	}
	err = json.Unmarshal(b, &fEntry)
	fEntry.FMRecordID = rec.ID
	if err != nil {
		return warehouse.Entry{}, api.NewError(err, fmt.Sprintf("没有找到id为 %s 的入库", id), "原因无知，请联系管理员")
	}

	return fEntry.ToEntry(), nil
}

func (s *EntryStore) GetEntryList(meta api.RequestMeta) ([]warehouse.Entry, api.ResponseMeta, error) {
	var resMeta api.ResponseMeta
	q := fm.NewFMQuery(s.databaseName, ENTRY_LAYOUT, fm.Find)
	var qFields []fm.FMQueryField
	for k, v := range meta.InternalFilter {
		qFields = append(qFields, fm.FMQueryField{Name: k, Value: v, Op: "="})
	}
	q.WithFields(qFields...)
	recs, resMeta, err := fmutil.GetFileMakerRecordList(s, q, meta)
	if err != nil {
		return nil, resMeta, api.NewError(err, "无法获取入库列表", "原因无知，请联系管理员")
	}

	var fEntries []warehouse.FileMakerEntry
	for _, rec := range recs {
		entry := warehouse.FileMakerEntry{}
		b, err := rec.JsonFields()
		if err != nil {
			return nil, resMeta, api.NewError(err, "无法获取入库列表", "原因无知，请联系管理员")
		}

		err = json.Unmarshal(b, &entry)
		if err != nil {
			return nil, resMeta, api.NewError(err, "无法获取入库列表", "原因无知，请联系管理员")
		}
		entry.FMRecordID = rec.ID
		fEntries = append(fEntries, entry)
	}

	var entries []warehouse.Entry
	for _, fe := range fEntries {
		entries = append(entries, fe.ToEntry())
	}

	return entries, resMeta, nil
}

func (s *EntryStore) CreateEntry(e warehouse.Entry) (warehouse.Entry, error) {
	q := fm.NewFMQuery(s.databaseName, ENTRY_LAYOUT, fm.New)
	q.WithFields(
		fm.FMQueryField{Name: "CustomerCode", Value: e.CustomerCode},
		fm.FMQueryField{Name: "SourceOfEntry", Value: e.Source},
		fm.FMQueryField{Name: "TrackCode", Value: e.TrackCode},
		fm.FMQueryField{Name: "QuantityOfBoxes", Value: strconv.Itoa(e.BoxQty)},
		fm.FMQueryField{Name: "PieceQuantity", Value: strconv.Itoa(e.PcsQty)},
		fm.FMQueryField{Name: "ProductName", Value: e.ProductName},
		fm.FMQueryField{Name: "Warehouse", Value: e.Warehouse},
		fm.FMQueryField{Name: "is_found_for_shipment", Value: strconv.Itoa(fmutil.ConvertBoolToInt(e.IsFoundForShipment))},
		fm.FMQueryField{Name: "has_brand", Value: strconv.Itoa(fmutil.ConvertBoolToInt(e.HasBrand))},
		fm.FMQueryField{Name: "product_category", Value: string(e.ProductCategory)},
		fm.FMQueryField{Name: "CreatedBy_Account", Value: s.conn.Username},
	)

	fmSet, err := s.conn.Query(q)
	if err != nil {
		return warehouse.Entry{}, api.NewError(err, "入库创建失败", "原因无知，请联系管理员")
	}

	resEntry := warehouse.FileMakerEntry{}
	if len(fmSet.Resultset.Records) < 1 {
		return warehouse.Entry{}, api.NewError(ErrZeroRecordsInResultSet, "入库创建失败", "原因无知，请联系管理员")
	}

	b, err := fmSet.Resultset.Records[0].JsonFields()
	if err != nil {
		return warehouse.Entry{}, api.NewError(err, "入库创建失败", "原因无知，请联系管理员")
	}

	err = json.Unmarshal(b, &resEntry)
	if err != nil {
		return warehouse.Entry{}, api.NewError(err, "入库创建失败", "原因无知，请联系管理员")
	}

	return resEntry.ToEntry(), nil
}

func (s *EntryStore) UpdateEntry(e warehouse.Entry) (*warehouse.Entry, error) {
	q := fm.NewFMQuery(s.databaseName, ENTRY_LAYOUT, fm.Edit)
	q.WithRecordId(e.FMRecordID)
	q.WithFields(
		fm.FMQueryField{Name: "CustomerCode", Value: e.CustomerCode},
		fm.FMQueryField{Name: "SourceOfEntry", Value: e.Source},
		fm.FMQueryField{Name: "TrackCode", Value: e.TrackCode},
		fm.FMQueryField{Name: "QuantityOfBoxes", Value: strconv.Itoa(e.BoxQty)},
		fm.FMQueryField{Name: "PieceQuantity", Value: strconv.Itoa(e.PcsQty)},
		fm.FMQueryField{Name: "ProductName", Value: e.ProductName},
		fm.FMQueryField{Name: "Warehouse", Value: e.Warehouse},
		fm.FMQueryField{Name: "is_found_for_shipment", Value: strconv.Itoa(fmutil.ConvertBoolToInt(e.IsFoundForShipment))},
		fm.FMQueryField{Name: "has_brand", Value: strconv.Itoa(fmutil.ConvertBoolToInt(e.HasBrand))},
		fm.FMQueryField{Name: "product_category", Value: string(e.ProductCategory)},
	)

	var auditData string
	for _, fg := range q.QueryFields {
		for _, f := range fg.Fields {
			auditData += fmt.Sprintf("[%s:%s]", f.Name, f.Value)
		}
	}

	b, err := json.Marshal(e)
	if err != nil {
		return nil, api.NewError(err, fmt.Sprintf("更新入库 %s 失败", e.ID), "原因无知，请联系管理员")
	}
	fmutil.WithAudit(q, e.ID, "Entries", "api_edit_record", auditData, s.conn.Username)

	fmSet, err := s.conn.Query(q)
	if err != nil {
		return nil, api.NewError(err, fmt.Sprintf("更新入库 %s 失败", e.ID), "原因无知，请联系管理员")
	}

	resEntry := warehouse.FileMakerEntry{}
	if len(fmSet.Resultset.Records) < 1 {
		return nil, api.NewError(errors.New("0 records in fmResultSet after entry update"), fmt.Sprintf("更新入库 %s 失败", e.ID), "原因无知，请联系管理员")
	}

	b, err = fmSet.Resultset.Records[0].JsonFields()
	if err != nil {
		return nil, api.NewError(err, fmt.Sprintf("更新入库 %s 失败", e.ID), "原因无知，请联系管理员")
	}

	err = json.Unmarshal(b, &resEntry)
	if err != nil {
		return nil, api.NewError(err, fmt.Sprintf("更新入库 %s 失败", e.ID), "原因无知，请联系管理员")
	}

	updatedEntry := resEntry.ToEntry()

	return &updatedEntry, nil
}
