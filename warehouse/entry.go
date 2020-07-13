//go:generate gomodifytags -file $GOFILE -struct Entry -add-tags json -w  -add-options json=omitempty

package warehouse

import (
	"github.com/amanbolat/ca-warehouse-client/api"
	"github.com/amanbolat/ca-warehouse-client/filemaker/fmutil"
	"time"
)

var entryFieldNamesMap = map[string]string{
	"id":              "id",
	"customer_code":   "CustomerCode",
	"shipment_code":   "Id_shipmentNumber",
	"date_of_entry":   "Date_Created_Timestamp",
	"source_of_entry": "SourceOfEntry",
	"track_code":      "TrackCode",
	"box_qty":         "QuantityOfBoxes",
	"pcs_qty":         "PieceQuantity",
	"product_name":    "ProductName",
	"warehouse":       "Warehouse",
}

// MapEntryFields maps api field names into FileMaker field names
func MapEntryFields(meta api.RequestMeta) api.RequestMeta {
	var newMeta api.RequestMeta
	newMeta = meta
	newMeta.SortFields = []api.SortField{}
	newMeta.InternalFilter = make(map[string]string)

	for _, field := range meta.SortFields {
		n, ok := entryFieldNamesMap[field.Name]
		if ok {
			newMeta.SortFields = append(newMeta.SortFields, api.SortField{
				Name:       n,
				Descending: field.Descending,
			})
		}
	}

	for _, filter := range meta.Filters {
		key, ok := entryFieldNamesMap[filter.K]
		if ok {
			newMeta.InternalFilter[key] = filter.V
		}
	}

	return newMeta
}

type EntryStatus string

const (
	EntryStatusUtilized EntryStatus = "utilized"
	EntryStatusReceived EntryStatus = "received"
	EntryStatusPacked   EntryStatus = "packed"
	EntryStatusSentOut  EntryStatus = "sent_out"
)

type Entry struct {
	ID                 string          `json:"id,omitempty"`
	CustomerCode       string          `json:"customer_code,omitempty"`
	ShipmentCode       string          `json:"shipment_code,omitempty"`
	Status             int             `json:"status,omitempty"`
	DateOfEntry        time.Time       `json:"date_of_entry,omitempty"`
	Source             string          `json:"source_of_entry,omitempty"`
	TrackCode          string          `json:"track_code,omitempty"`
	BoxQty             int             `json:"box_qty,omitempty"`
	PcsQty             int             `json:"pcs_qty,omitempty"`
	ProductName        string          `json:"product_name,omitempty"`
	Warehouse          string          `json:"warehouse,omitempty"`
	ImageUrls          []string        `json:"image_urls,omitempty"`
	HasBrand           bool            `json:"has_brand"`
	IsFoundForShipment bool            `json:"is_found_for_shipment"`
	ProductCategory    ProductCategory `json:"product_category,omitempty"`
	FMRecordID         int             `json:"-,omitempty"`
}

type FileMakerEntry struct {
	ID                 string    `json:"id"`
	CustomerCode       string    `json:"CustomerCode"`
	ShipmentNumber     string    `json:"Id_shipmentNumber,omitempty"`
	Status             int       `json:"StatusOfEntry_key"`
	DateOfEntry        time.Time `json:"Date_Created_Timestamp,omitempty"`
	Source             string    `json:"SourceOfEntry,omitempty"`
	TrackCode          string    `json:"TrackCode,omitempty"`
	BoxQty             float64   `json:"QuantityOfBoxes,omitempty"`
	PcsQty             float64   `json:"PieceQuantity,omitempty"`
	ProductName        string    `json:"ProductName,omitempty"`
	Warehouse          string    `json:"Warehouse,omitempty"`
	ImageUrls          []string  `json:"Container,omitempty"`
	IsFoundForShipment int       `json:"is_found_for_shipment"`
	HasBrand           int       `json:"has_brand"`
	ProductCategory    string    `json:"product_category"`
	ShipmentStatusKey  int       `json:"TO4a_Entries||Shipments::ShipmentStatus_number"`
	FMRecordID         int       `json:"-"`
}

func (v *FileMakerEntry) ToEntry() Entry {
	return Entry{
		ID:                 v.ID,
		CustomerCode:       v.CustomerCode,
		ShipmentCode:       v.ShipmentNumber,
		Status:             v.Status,
		DateOfEntry:        v.DateOfEntry,
		Source:             v.Source,
		TrackCode:          v.TrackCode,
		BoxQty:             int(v.BoxQty),
		PcsQty:             int(v.PcsQty),
		ProductName:        v.ProductName,
		Warehouse:          v.Warehouse,
		ImageUrls:          v.ImageUrls,
		IsFoundForShipment: fmutil.ConvertToBool(v.IsFoundForShipment),
		HasBrand:           fmutil.ConvertToBool(v.HasBrand),
		ProductCategory:    ProductCategory(v.ProductCategory),
		FMRecordID:         v.FMRecordID,
	}
}
