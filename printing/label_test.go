package printing

import (
	"github.com/amanbolat/ca-warehouse-client/crm"
	"github.com/amanbolat/ca-warehouse-client/logistics"
	"github.com/amanbolat/ca-warehouse-client/warehouse"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

var fontPath = os.Getenv("FONT_PATH")
var sp = logistics.Shipment{
	Code:         "SPN007001",
	CustomerCode: "77-00123",
	Entries: []*warehouse.Entry{
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 12, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 10, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 44, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 33, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 55, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
		{ID: "EN0001", Source: "顺丰快递", TrackCode: "SF1241923123", BoxQty: 1, PcsQty: 100},
	},
	PackageMethod:   "wooden_crate",
	PackageMethodZh: "泡沫",
	UnitLoads: []*logistics.UnitLoad{
		{
			Sequence:    1,
			Quantity:    100,
			Weight:      decimal.NewFromFloat(5555.55),
			Length:      100,
			Height:      33,
			Width:       55,
			ProductName: "灯具，LED 包包 衣服 灯具，LED 包包 衣服",
		},
		{
			Sequence:    2,
			Quantity:    111,
			Weight:      decimal.NewFromFloat(100.55),
			Length:      100,
			Height:      44,
			Width:       67,
			ProductName: "LED 包包 衣服",
		},
		{
			Sequence:    3,
			Quantity:    111,
			Weight:      decimal.NewFromFloat(77.10),
			Length:      89,
			Height:      44,
			Width:       67,
			ProductName: "灯具，LED 包包 衣服 裤子",
		},
	},
	NeedDeclare: true,
	PartnerInfo: logistics.PartnerInfo{
		Code:           "XX-ZHANGYUNLONG",
		DeliveryMethod: logistics.DMLandContainer,
		ProductName:    "LED灯具100个，机子200个，电池300个，配件100个，轮子300，LEDDDDDDDD灯具100个，机子200个，电池300个，配件100个，轮子300",
		CargoValue:     0,
		Recipient: logistics.Recipient{
			Name:        "Михаил задорный 波哥",
			PhoneNumber: "7800123123123",
			Destination: "Москва 莫斯科",
		},
	},
	Notes: []*crm.Note{
		{Content: "加护角"},
		{Content: "可以发了"},
		{Content: "尽量往木箱里添加泡沫+纸板什么的"},
	},
}

func TestBarcodeManger_CreateEntryBarcode(t *testing.T) {
	entryID := "en00001231"
	fontPath := os.Getenv("FONT_PATH")

	bm, err := NewLabelManger(fontPath)
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
	bm, err := NewLabelManger(fontPath)
	assert.NoError(t, err)
	labels, err := bm.CreateUnitLoadLabels(sp)
	if assert.NoError(t, err) {
		exec.Command("open", labels.FullPath).Run()
	}
}

func TestLabelManager_CreateShipmentPreparationLabel(t *testing.T) {
	lm, err := NewLabelManger(fontPath)
	assert.NoError(t, err)

	entriesLabel, err := lm.CreateShipmentPreparationLabels(sp)
	assert.NoError(t, err)
	if assert.NoError(t, err) {
		exec.Command("open", entriesLabel.FullPath).Run()
	}
}

func TestLabelManager_CreateShipmentPartnerInfoLabel(t *testing.T) {
	lm, err := NewLabelManger(fontPath)
	assert.NoError(t, err)

	partnerInfoLabel, err := lm.CreateShipmentPartnerInfoLabel(sp)
	assert.NoError(t, err)
	if assert.NoError(t, err) {
		exec.Command("open", partnerInfoLabel.FullPath).Run()
	}
}
