package filemaker

import (
	"encoding/json"
	query "github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/filemaker/fmutil"
	"github.com/amanbolat/ca-warehouse-client/logistics"
	fm "github.com/amanbolat/gofmcon"
	"github.com/pkg/errors"
)

const (
	SHIPMENT_LAYOUT = "warehouse_shipment_single"
)

type ShipmentStore struct {
	conn         *fm.FMConnector
	databaseName string
}

func (r *ShipmentStore) DBName() string {
	return r.databaseName
}

func (r *ShipmentStore) FMConn() *fm.FMConnector {
	return r.conn
}

func NewShipmentStore(conn *fm.FMConnector, dbName string) *ShipmentStore {
	return &ShipmentStore{conn, dbName}
}

// updateShipment updates shipment
func (r *ShipmentStore) updateShipment(q *fm.FMQuery) (logistics.Shipment, error) {
	fmSet, err := r.conn.Query(q)
	if err != nil {
		return logistics.Shipment{}, errors.WithMessage(err, "database_error")
	}

	if len(fmSet.Resultset.Records) < 1 {
		return logistics.Shipment{}, errors.New("database_error: update failed")
	}

	fShipment := logistics.FileMakerShipment{}
	b, err := fmSet.Resultset.Records[0].JsonFields()
	if err != nil {
		return logistics.Shipment{}, err
	}
	err = json.Unmarshal(b, &fShipment)
	if err != nil {
		return logistics.Shipment{}, err
	}

	return fShipment.ToShipment(), nil
}

func (r *ShipmentStore) GetShipmentList(meta query.RequestMeta) ([]*logistics.Shipment, query.ResponseMeta, error) {
	var resMeta query.ResponseMeta
	q := fm.NewFMQuery(r.databaseName, SHIPMENT_LAYOUT, fm.Find)
	q.WithFields(
		fm.FMQueryField{
			Name:  "Departure_Warehouse",
			Value: "GZWH2",
			Op:    fm.Equal,
		},
		fm.FMQueryField{
			Name:  "ShipmentStatus_number",
			Value: "1...2",
			Op:    "=",
		})

	recs, resMeta, err := fmutil.GetFileMakerRecordList(r, q, meta)
	if err != nil {
		return nil, resMeta, err
	}

	var shipments []*logistics.Shipment
	for _, rec := range recs {
		fShipment := logistics.FileMakerShipment{}
		b, err := rec.JsonFields()
		if err != nil {
			return nil, resMeta, err
		}
		err = json.Unmarshal(b, &fShipment)
		if err != nil {
			return nil, resMeta, err
		}
		s := fShipment.ToShipment()
		shipments = append(shipments, &s)
	}

	return shipments, resMeta, nil
}

func (r *ShipmentStore) GetShipmentByCode(code string) (logistics.Shipment, error) {
	q := fm.NewFMQuery(r.databaseName, SHIPMENT_LAYOUT, fm.Find)
	q.WithResponseLayout(SHIPMENT_LAYOUT)
	q.WithFields(
		fm.FMQueryField{Name: "code", Value: code, Op: fm.Equal},
	)

	rec, err := fmutil.GetFileMakerRecordSingle(r, q)
	if err != nil {
		return logistics.Shipment{}, err
	}

	fShipment := logistics.FileMakerShipment{}
	b, err := rec.JsonFields()
	if err != nil {
		return logistics.Shipment{}, err
	}
	err = json.Unmarshal(b, &fShipment)
	if err != nil {
		return logistics.Shipment{}, err
	}

	return fShipment.ToShipment(), nil
}
