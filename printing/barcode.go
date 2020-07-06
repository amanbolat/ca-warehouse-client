package printing

import (
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/logistics"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fogleman/gg"
	"github.com/rs/xid"
	"gopkg.in/errgo.v2/errors"
	"os"
	"path"
	"strings"
)

const PAPER_W, PAPER_H = 294, 447

type BarcodeManger struct {
	fontPath string
}

//
func NewBarcodeManger(fontPath string) (*BarcodeManger, error) {
	info, err := os.Lstat(fontPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("fonts are not found in: %s", fontPath))
	}
	if info.IsDir() {
		return nil, errors.New("provided font path is not file")
	}
	return &BarcodeManger{fontPath: fontPath}, nil
}

type Barcode struct {
	File     *os.File
	FullPath string
}

func (bm BarcodeManger) CreateEntryBarcode(entryId string) (Barcode, error) {
	// text is lowered, because bcst-50 scanner has problem with scanning
	// UpperCase chars
	entryId = strings.ToLower(entryId)
	res := Barcode{}
	tmpFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-entry-barcode.png", xid.New()))

	file, err := os.Create(tmpFilePath)
	if err != nil {
		return res, err
	}

	defer file.Close()

	br, err := qr.Encode(entryId, qr.M, qr.Auto)
	if err != nil {
		return res, err
	}

	bcWidth, bcHeight := 200, 200

	bc, err := barcode.Scale(br, bcWidth, bcHeight)
	if err != nil {
		return res, err
	}

	dc := gg.NewContext(PAPER_W, PAPER_H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	err = dc.LoadFontFace(bm.fontPath, 33)
	if err != nil {
		return res, err
	}
	dc.DrawImage(bc, (PAPER_W-bcWidth)/2, (PAPER_H-bcHeight)/2)
	dc.DrawStringAnchored(strings.ToUpper(entryId), float64(PAPER_W/2), float64((PAPER_H-bcHeight)/2+bcHeight+20), 0.5, 0.5)
	err = dc.EncodePNG(file)

	if err != nil {
		return res, err
	}

	res.File = file
	res.FullPath = tmpFilePath

	return res, nil
}

func (bm BarcodeManger) CreateUnitLoadBarcodes(shipment logistics.Shipment) ([]Barcode, error) {
	var barcodes []Barcode
	ulCount := len(shipment.UnitLoads)

	for _, ul := range shipment.UnitLoads {
		var ulBarcode Barcode
		tmpFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-unitload-barcode.png", xid.New()))
		file, err := os.Create(tmpFilePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		dc := gg.NewContext(PAPER_W, PAPER_H)
		dc.SetRGB(1, 1, 1)
		dc.Clear()
		dc.SetRGB(0, 0, 0)
		err = dc.LoadFontFace(bm.fontPath, 20)
		if err != nil {
			return nil, err
		}

		bcWidth, bcHeight := 50, 50

		br, err := qr.Encode(fmt.Sprintf("%s-%d/%d", strings.ToLower(shipment.Code), ul.Sequence, ulCount), qr.H, qr.Auto)
		if err != nil {
			return nil, err
		}

		bc, err := barcode.Scale(br, bcWidth, bcHeight)
		if err != nil {
			return nil, err
		}

		dc.DrawImage(bc, 10, 10)
		dc.DrawStringAnchored(fmt.Sprintf("%s-%d/%d", strings.ToUpper(shipment.Code), ul.Sequence, ulCount), 70, 45, 0, 0)
		dc.SetLineWidth(1)
		dc.DrawLine(10, 70, PAPER_W-10, 70)
		dc.Stroke()
		dc.LoadFontFace(bm.fontPath, 16)
		ulInfo := fmt.Sprintf("Габариты места/此包尺寸")
		dc.DrawStringAnchored(ulInfo, float64(PAPER_W/2), 95, 0.5, 0.5)
		ulSize := fmt.Sprintf("%s Kg    %s m3    %d × %d × %d", ul.Weight, ul.Cubage(), ul.Height, ul.Length, ul.Width)
		dc.DrawStringAnchored(ulSize, float64(PAPER_W/2), 130, 0.5, 0.5)
		dc.DrawLine(10, 165, PAPER_W-10, 165)
		dc.Stroke()
		shipmentInfo := fmt.Sprintf("Общие габариты/综合尺寸")
		dc.DrawStringAnchored(shipmentInfo, float64(PAPER_W/2), 190, 0.5, 0.5)
		sWeight := fmt.Sprintf("%s Kg    %s m3", shipment.Weight(), shipment.Cubage())
		dc.DrawStringAnchored(sWeight, float64(PAPER_W/2), 225, 0.5, 0.5)

		// Attention text
		dc.DrawRectangle(10, 260, PAPER_W-20, 75)
		dc.Stroke()
		attentionText := `Внимание! При получении груза обязательно проверьте целостность внешней упаковки и характеристики груза`
		dc.DrawStringWrapped(attentionText, 20, 270, 0, 0, PAPER_W-30, 1.2, gg.AlignLeft)

		// Contacts
		dc.LoadFontFace(bm.fontPath, 12)
		dc.DrawString("Кросс Азия", 10, 380)
		dc.DrawString("Телефон 1: +7(812)309-73-97", 10, 395)
		dc.DrawString("Телефон 2: +86 136-9923-3755", 10, 410)
		dc.DrawString("Сайт: https://crossasia.ru", 10, 425)

		err = dc.EncodePNG(file)

		if err != nil {
			return nil, err
		}

		ulBarcode.File = file
		ulBarcode.FullPath = tmpFilePath

		barcodes = append(barcodes, ulBarcode)
	}

	return barcodes, nil
}
