//go:generate gomodifytags -file $GOFILE -struct Entry -add-tags json -add-options json=omitempty -w

package warehouse

import (
	"time"
)

type Entry struct {
	ID           string    `json:"id,omitempty"`
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
	FMRecordID   int64     `json:"-"`
}

type FileMakerEntry struct {
	ID             string    `json:"id"`
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
