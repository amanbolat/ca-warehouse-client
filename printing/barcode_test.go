package printing

import (
	"github.com/amanbolat/ca-warehouse-client/logistics"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestBarcodeManger_CreateEntryBarcode(t *testing.T) {
	entryID := "en00001231"
	fontPath := os.Getenv("FONT_PATH")

	bm, err := NewBarcodeManger(fontPath)
	assert.NoError(t, err)
	bc, err := bm.CreateEntryBarcode(entryID)
	assert.NoError(t, err)
	assert.NotNil(t, bc.File)
	assert.FileExists(t, bc.FullPath)
	if assert.FileExists(t, bc.FullPath) {
		exec.Command("open", bc.FullPath).Run()
	}
}

func TestBarcodeManger_CreateUnitLoadBarcodes(t *testing.T) {
	fontPath := os.Getenv("FONT_PATH")
	sp := logistics.Shipment{
		Code: "SPN007001",
		UnitLoads: []*logistics.UnitLoad{
			{
				Sequence: 1,
				Quantity: 100,
				Weight:   decimal.NewFromFloat(50.55),
				Length:   100,
				Height:   33,
				Width:    55,
			},
		},
	}

	bm, err := NewBarcodeManger(fontPath)
	assert.NoError(t, err)
	barcodes, err := bm.CreateUnitLoadBarcodes(sp)

	for _, bc := range barcodes {
		assert.NoError(t, err)
		assert.NotNil(t, bc.File)
		assert.FileExists(t, bc.FullPath)
		if assert.FileExists(t, bc.FullPath) {
			exec.Command("open", bc.FullPath).Run()
		}
	}
}
