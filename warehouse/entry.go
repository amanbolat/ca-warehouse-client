//go:generate gomodifytags -file $GOFILE -struct Entry -add-tags json -add-options json=omitempty -w

package warehouse

import (
	"github.com/amanbolat/ca-warehouse-client/api"
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

	for _, field := range meta.SortFields {
		newMeta.SortFields = append(newMeta.SortFields, api.SortField{
			Name:       entryFieldNamesMap[field.Name],
			Descending: field.Descending,
		})
	}

	return newMeta
}

type Entry struct {
	ID           string    `json:"id,omitempty"`
	CustomerCode string    `json:"customer_code,omitempty"`
	ShipmentCode string    `json:"shipment_code,omitempty"`
	Status       int       `json:"status,omitempty"`
	DateOfEntry  time.Time `json:"date_of_entry,omitempty"`
	Source       string    `json:"source_of_entry,omitempty"`
	TrackCode    string    `json:"track_code,omitempty"`
	BoxQty       int       `json:"box_qty,omitempty"`
	PcsQty       int       `json:"pcs_qty,omitempty"`
	ProductName  string    `json:"product_name,omitempty"`
	Warehouse    string    `json:"warehouse,omitempty"`
	ImageUrls    []string  `json:"image_urls,omitempty"`
	FMRecordID   int64     `json:"-,omitempty"`
}

type FileMakerEntry struct {
	ID             string    `json:"id"`
	CustomerCode   string    `json:"CustomerCode"`
	ShipmentNumber string    `json:"Id_shipmentNumber,omitempty"`
	Status         int       `json:"StatusOfEntry_key"`
	DateOfEntry    time.Time `json:"Date_Created_Timestamp,omitempty"`
	Source         string    `json:"SourceOfEntry,omitempty"`
	TrackCode      string    `json:"TrackCode,omitempty"`
	BoxQty         float64   `json:"QuantityOfBoxes,omitempty"`
	PcsQty         float64   `json:"PieceQuantity,omitempty"`
	ProductName    string    `json:"ProductName,omitempty"`
	Warehouse      string    `json:"Warehouse,omitempty"`
	ImageUrls      []string  `json:"Container,omitempty"`
	FMRecordID     int64     `json:"-"`
}

func (fe FileMakerEntry) ToEntry() Entry {
	return Entry{
		ID:           fe.ID,
		CustomerCode: fe.CustomerCode,
		ShipmentCode: fe.ShipmentNumber,
		Status:       fe.Status,
		DateOfEntry:  fe.DateOfEntry,
		Source:       fe.Source,
		TrackCode:    fe.TrackCode,
		BoxQty:       int(fe.BoxQty),
		PcsQty:       int(fe.PcsQty),
		ProductName:  fe.ProductName,
		Warehouse:    fe.Warehouse,
		ImageUrls:    fe.ImageUrls,
		FMRecordID:   fe.FMRecordID,
	}
}
