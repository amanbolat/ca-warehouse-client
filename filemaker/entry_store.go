package filemaker

import (
	"encoding/json"
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	fm "github.com/amanbolat/gofmcon"
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

	rec, err := GetFileMakerRecordSingle(s, q)

	if err != nil {
		return warehouse.Entry{}, err
	}

	fEntry := warehouse.FileMakerEntry{}
	b, err := rec.JsonFields()
	if err != nil {
		return warehouse.Entry{}, err
	}
	err = json.Unmarshal(b, &fEntry)
	if err != nil {
		return warehouse.Entry{}, err
	}

	return fEntry.ToEntry(), nil
}

func (s *EntryStore) GetEntryList(meta api.RequestMeta) ([]warehouse.Entry, api.ResponseMeta, error) {
	var resMeta api.ResponseMeta
	q := fm.NewFMQuery(s.databaseName, ENTRY_LAYOUT, fm.Find)
	q.WithFields(
		fm.FMQueryField{
			Name:  "Warehouse",
			Value: "GZWH1",
			Op:    fm.Equal,
		},
		fm.FMQueryField{
			Name:  "Id_shipmentNumber",
			Value: "",
			Op:    fm.Equal,
		},
		fm.FMQueryField{
			Name:  "is_utilized",
			Value: "",
			Op:    fm.Equal,
		},
	)
	recs, resMeta, err := GetFileMakerRecordList(s, q, meta)
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
		fEntries = append(fEntries, entry)
	}

	var entries []warehouse.Entry
	for _, fe := range fEntries {
		entries = append(entries, fe.ToEntry())
	}

	return entries, resMeta, err
}
