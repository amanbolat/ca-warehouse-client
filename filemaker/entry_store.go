package filemaker

import (
	"encoding/json"
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/filemaker/fmutil"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	fm "github.com/amanbolat/gofmcon"
	"github.com/pkg/errors"
	"strconv"
)

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
		return warehouse.Entry{}, err
	}

	fEntry := warehouse.FileMakerEntry{}
	b, err := rec.JsonFields()
	if err != nil {
		return warehouse.Entry{}, err
	}
	err = json.Unmarshal(b, &fEntry)
	fEntry.FMRecordID = rec.ID
	if err != nil {
		return warehouse.Entry{}, err
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
		return nil, resMeta, err
	}

	var fEntries []warehouse.FileMakerEntry
	for _, rec := range recs {
		entry := warehouse.FileMakerEntry{}
		b, err := rec.JsonFields()
		if err != nil {
			return nil, resMeta, err
		}

		err = json.Unmarshal(b, &entry)
		if err != nil {
			return nil, resMeta, err
		}
		entry.FMRecordID = rec.ID
		fEntries = append(fEntries, entry)
	}

	var entries []warehouse.Entry
	for _, fe := range fEntries {
		entries = append(entries, fe.ToEntry())
	}

	return entries, resMeta, err
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
	)

	fmSet, err := s.conn.Query(q)
	if err != nil {
		return warehouse.Entry{}, errors.WithMessage(err, "could not create new entry")
	}

	resEntry := warehouse.FileMakerEntry{}
	if len(fmSet.Resultset.Records) < 1 {
		return warehouse.Entry{}, errors.New("new entry might be created, but no result was returned from database")
	}

	b, err := fmSet.Resultset.Records[0].JsonFields()
	if err != nil {
		return warehouse.Entry{}, err
	}

	err = json.Unmarshal(b, &resEntry)
	if err != nil {
		return warehouse.Entry{}, err
	}

	return resEntry.ToEntry(), nil
}
