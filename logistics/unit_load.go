//go:generate gomodifytags -file $GOFILE -struct UnitLoad -add-tags json -w
//go:generate gomodifytags -file $GOFILE -struct jsonUnitLoad -add-tags json -w
package logistics

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type UnitLoad struct {
	Sequence    int             `json:"sequence"`
	Quantity    int             `json:"quantity"`
	ProductName string          `json:"product_name"`
	Weight      decimal.Decimal `json:"weight"` // kg
	Length      int64           `json:"length"` // cm
	Height      int64           `json:"height"` // cm
	Width       int64           `json:"width"`  // cm
	FMRecordID  int64           `json:"-"`
}

type FileMakerUnitLoad struct {
	Sequence    int     `json:"SequenceNumber"`
	Quantity    int     `json:"Quantity"`
	ProductName string  `json:"SD_ProductName"`
	Weight      float64 `json:"SD_Weight"`
	Length      int64   `json:"SD_Length"` // cm
	Height      int64   `json:"SD_Height"` // cm
	Width       int64   `json:"SD_Width"`  // cm
	FMRecordID  int64   `json:"-"`
}

func (fu FileMakerUnitLoad) ToUnitLoad() UnitLoad {
	return UnitLoad{
		Sequence:    fu.Sequence,
		Quantity:    fu.Quantity,
		ProductName: fu.ProductName,
		Weight:      decimal.NewFromFloat(fu.Weight),
		Length:      fu.Length,
		Height:      fu.Height,
		Width:       fu.Width,
		FMRecordID:  fu.FMRecordID,
	}
}

type AliasUnitLoad UnitLoad
type jsonUnitLoad struct {
	Cubage decimal.Decimal `json:"cubage"`
	*AliasUnitLoad
}

func (ul *UnitLoad) MarshalJSON() ([]byte, error) {
	pUL := &jsonUnitLoad{
		Cubage:        ul.Cubage(),
		AliasUnitLoad: (*AliasUnitLoad)(ul),
	}

	return json.Marshal(pUL)
}

// Cubage returns unit load volume in m3
func (ul *UnitLoad) Cubage() decimal.Decimal {
	if ul.Length == 0 || ul.Height == 0 || ul.Width == 0 {
		return decimal.Zero
	}

	cmCubage := ul.Length * ul.Height * ul.Width

	return decimal.New(cmCubage, 0).DivRound(decimal.New(1000000, 0), 2)
}
