//go:generate enumer -type=ShipmentStatus -json -sql -transform=snake
//go:generate enumer -type=ShipmentType -json -sql -transform=snake
//go:generate enumer -type=TransportMethod -json -sql -transform=snake
//go:generate enumer -type=PackageMethod -json -sql -transform=snake
//go:generate gomodifytags -file $GOFILE -struct Shipment -add-tags json -w
//go:generate gomodifytags -file $GOFILE -struct jsonShipment -add-tags json -w

package logistics

import (
	"encoding/json"
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

type ShipmentStatus int

const (
	Planning ShipmentStatus = iota
	Preparation
	Packed
	SentOut
	CustomsClearance
	OnTheWayToTP
	DeliveredToTP
	DeliveredToRecipient
	InvalidStatus ShipmentStatus = 999
)

type ShipmentType int

func ShipmentTypeP(i int) *ShipmentType {
	v := ShipmentType(i)
	return &v
}

const (
	CommonShipment ShipmentType = iota
	ConsolidationShipment
)

type TransportMethod int

func TransportMethodP(i int) *TransportMethod {
	v := TransportMethod(i)
	return &v
}

const (
	Air TransportMethod = iota
	Auto
	Train
	Express
	Sea
	Local
)

type PackageMethod int

func PackageMethodP(i int) *PackageMethod {
	v := PackageMethod(i)
	return &v
}

const (
	PackageNone PackageMethod = iota
	PackageBag
	PackageStandard
	PackageCarton
	PackageFoam
	PackageCartonFoam
	PackageWoodenCrate
	PackageWoodenCrateFoam
	PackageWoodenBox
	PackageWoodenBoxFoam
)

type FileMakerShipment struct {
	ID                     string    `json:"Id_shipment,omitempty"`
	Code                   string    `json:"code,omitempty"`
	Type                   int       `json:"CargoType_number,omitempty"`
	CustomerCode           string    `json:"CustomerCode,omitempty"`
	PackagesQty            int       `json:"PackageQuantity,omitempty"`
	PiecesQty              int       `json:"TotalQuanity,omitempty"`
	CurrentStatusKey       int       `json:"ShipmentStatus_number,omitempty"`
	TransferPointKey       int       `json:"TransferPoint_number,omitempty"`
	TransportMethodKey     int       `json:"TransportationMethod_number,omitempty"`
	PackageMethodKey       int       `json:"PackageMethod_number,omitempty"`
	DepartureWarehouse     string    `json:"Departure_Warehouse,omitempty"`
	ArrivalWarehouse       string    `json:"Arrival_Warehouse,omitempty"`
	TransferPointWarehouse string    `json:"TransferPoint_Warehouse,omitempty"`
	DateCreated            time.Time `json:"Date_Created,omitempty"`
	DateModified           time.Time `json:"Date_Modified_Timestamp,omitempty"`
	// Order values
	UnitLoads     []*FileMakerUnitLoad        `json:"TO2b_Shipments||ShipmentDetails,omitempty"`
	Entries       []*warehouse.FileMakerEntry `json:"TO2c_Shipments||Entries,omitempty"`
	Consolidation []*FileMakerShipment        `json:"TO2k_Shipments||Shipments||Child,omitempty"`
	FMRecordID    int64                       `json:"-"`
	ImageUrls     []string                    `json:"Container,omitempty"`
}

func (fs FileMakerShipment) ToShipment() Shipment {
	transferPoint, ok := transferPointsMap[fs.TransferPointKey]
	if !ok {
		transferPoint = "unknown"
	}
	var unitLoads []*UnitLoad
	for _, ful := range fs.UnitLoads {
		ul := ful.ToUnitLoad()
		unitLoads = append(unitLoads, &ul)
	}

	var entries []*warehouse.Entry
	for _, fe := range fs.Entries {
		e := fe.ToEntry()
		entries = append(entries, &e)
	}

	var consolidatedShipments []*Shipment
	for _, fs := range fs.Consolidation {
		s := fs.ToShipment()
		consolidatedShipments = append(consolidatedShipments, &s)
	}

	return Shipment{
		ID:                     fs.ID,
		Code:                   fs.Code,
		Type:                   ShipmentTypeP(fs.Type),
		CustomerCode:           fs.CustomerCode,
		PackagesQty:            fs.PackagesQty,
		PiecesQty:              fs.PiecesQty,
		CurrentStatusKey:       ShipmentStatus(fs.CurrentStatusKey),
		TransferPoint:          transferPoint,
		TransportMethod:        TransportMethodP(fs.TransportMethodKey),
		PackageMethod:          PackageMethodP(fs.PackageMethodKey),
		DepartureWarehouse:     fs.DepartureWarehouse,
		ArrivalWarehouse:       fs.ArrivalWarehouse,
		TransferPointWarehouse: fs.TransferPointWarehouse,
		DateCreated:            fs.DateCreated,
		DateModified:           fs.DateModified,
		UnitLoads:              unitLoads,
		Entries:                entries,
		Consolidation:          consolidatedShipments,
		FMRecordID:             fs.FMRecordID,
		ImageUrls:              nil,
	}
}

type Shipment struct {
	ID                     string           `json:"id,omitempty"`
	Code                   string           `json:"code,omitempty"`
	Type                   *ShipmentType    `json:"type,omitempty"`
	CustomerCode           string           `json:"customer_code,omitempty"`
	PackagesQty            int              `json:"packages_qty,omitempty"`
	PiecesQty              int              `json:"pieces_qty,omitempty"`
	CurrentStatusKey       ShipmentStatus   `json:"current_status,omitempty"`
	TransferPoint          string           `json:"transfer_point,omitempty"`
	TransportMethod        *TransportMethod `json:"transport_method,omitempty"`
	PackageMethod          *PackageMethod   `json:"package_method,omitempty"`
	DepartureWarehouse     string           `json:"departure_warehouse,omitempty"`
	ArrivalWarehouse       string           `json:"arrival_warehouse,omitempty"`
	TransferPointWarehouse string           `json:"transfer_point_warehouse,omitempty"`
	DateCreated            time.Time        `json:"date_created,omitempty"`
	DateModified           time.Time        `json:"date_modified,omitempty"`
	// Order values
	UnitLoads     []*UnitLoad        `json:"unit_loads,omitempty"`
	Entries       []*warehouse.Entry `json:"entries,omitempty"`
	Consolidation []*Shipment        `json:"consolidation,omitempty"`
	FMRecordID    int64              `json:"-"`
	ImageUrls     []string           `json:"image_urls,omitempty"`
}

type AliasShipment Shipment
type jsonShipment struct {
	Weight  decimal.Decimal `json:"weight,omitempty"`
	Cubage  decimal.Decimal `json:"cubage,omitempty"`
	Density decimal.Decimal `json:"density,omitempty"`
	*AliasShipment
}

func (s Shipment) MarshalJSON() ([]byte, error) {
	ps := &jsonShipment{
		Weight:        s.Weight(),
		Cubage:        s.Cubage(),
		Density:       s.Density(),
		AliasShipment: (*AliasShipment)(&s),
	}

	return json.Marshal(ps)
}

func (s *Shipment) ToJSON() string {
	b, _ := json.Marshal(s)
	return string(b)
}

// Weight returns shipments weight in kg
func (s Shipment) Weight() decimal.Decimal {
	var w decimal.Decimal
	for _, e := range s.UnitLoads {
		w = w.Add(e.Weight)
	}

	return w.Round(2)
}

// Cubage returns shipment volume in m3
func (s Shipment) Cubage() decimal.Decimal {
	var c decimal.Decimal
	for _, e := range s.UnitLoads {
		c = c.Add(e.Cubage())
	}

	return c.Round(2)
}

func (s Shipment) Density() decimal.Decimal {
	if s.Cubage().Equal(decimal.Zero) {
		return decimal.Zero
	}
	return s.Weight().Div(s.Cubage()).Round(2)
}

func (s Shipment) CurrentStatus() string {
	return fmt.Sprintf("%s", s.CurrentStatusKey)
}

func (i ShipmentStatus) NextValid() ShipmentStatus {
	if i == DeliveredToRecipient {
		return InvalidStatus
	}

	return i + 1
}

func (s *Shipment) ChangeStatus(sts ShipmentStatus) error {
	if s.CurrentStatusKey.NextValid() != sts {
		return errors.New("shipment.ChangeStatus: invalid shipment status")
	}

	s.CurrentStatusKey = sts

	return nil
}

func (s *Shipment) AddUnitLoad(ul *UnitLoad) error {
	if s.CurrentStatusKey != Preparation {
		return errors.Errorf("shipment.AddUnitLoad: shipment should be on %s status", Preparation)
	}

	s.UnitLoads = append(s.UnitLoads, ul)

	return nil
}

func (s *Shipment) ResourceName() string {
	return "shipment"
}
