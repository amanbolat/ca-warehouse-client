package filemaker

import (
	"encoding/json"
	query "github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/crm"
	"github.com/amanbolat/ca-warehouse-client/filemaker/fmutil"
	fm "github.com/amanbolat/gofmcon"
)

const (
	CUSTOMER_LAYOUT = "warehouse_customer_list"
)

type CustomerStore struct {
	conn         *fm.FMConnector
	databaseName string
}

func (r *CustomerStore) DBName() string {
	return r.databaseName
}

func (r *CustomerStore) FMConn() *fm.FMConnector {
	return r.conn
}

func NewCustomerStore(conn *fm.FMConnector, dbName string) *CustomerStore {
	return &CustomerStore{conn, dbName}
}

func (r *CustomerStore) GetCustomerList(meta query.RequestMeta) ([]crm.Customer, query.ResponseMeta, error) {
	var resMeta query.ResponseMeta
	q := fm.NewFMQuery(r.databaseName, CUSTOMER_LAYOUT, fm.FindAll)

	recs, resMeta, err := fmutil.GetFileMakerRecordList(r, q, meta)
	if err != nil {
		return nil, resMeta, err
	}

	var customers []crm.Customer
	for _, rec := range recs {
		fCustomer := crm.FileMakerCustomer{}
		b, err := rec.JsonFields()
		if err != nil {
			return nil, resMeta, err
		}
		err = json.Unmarshal(b, &fCustomer)
		if err != nil {
			return nil, resMeta, err
		}
		c := fCustomer.ToCustomer()
		customers = append(customers, c)
	}

	return customers, resMeta, nil
}
