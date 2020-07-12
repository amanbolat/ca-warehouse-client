package printing

import (
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/i18n"
	"github.com/amanbolat/ca-warehouse-client/logistics"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/signintech/gopdf"
	"os"
	"path"
	"strings"
	"time"
)

const PAPER_W, PAPER_H = 294, 447

type LabelManager struct {
	fontPath string
}

func NewLabelManger(fontPath string) (LabelManager, error) {
	info, err := os.Lstat(fontPath)
	if err != nil {
		return LabelManager{}, errors.New(fmt.Sprintf("fonts are not found in: %s", fontPath))
	}
	if info.IsDir() {
		return LabelManager{}, errors.New("provided font path is not file")
	}
	return LabelManager{fontPath: fontPath}, nil
}

type Label struct {
	File     *os.File
	FullPath string
}

func (lm LabelManager) CreateEntryBarcode(entryId string) (Label, error) {
	// text is lowered, because bcst-50 scanner has problem with scanning
	// UpperCase chars
	entryId = strings.ToLower(entryId)
	res := Label{}
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
	err = dc.LoadFontFace(lm.fontPath, 33)
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

func (lm LabelManager) CreateUnitLoadLabels(shipment logistics.Shipment) ([]Label, error) {
	var barcodes []Label
	ulCount := len(shipment.UnitLoads)
	if ulCount < 1 {
		return nil, errors.New("there is now unit loads, nothing to print")
	}

	for _, ul := range shipment.UnitLoads {
		var ulBarcode Label
		tmpFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-unitload-barcode.png", xid.New()))
		file, err := os.Create(tmpFilePath)
		if err != nil {
			return nil, errors.WithMessage(err, "could not create tmp file")
		}
		defer file.Close()

		dc := gg.NewContext(PAPER_W, PAPER_H)
		dc.SetRGB(1, 1, 1)
		dc.Clear()
		dc.SetRGB(0, 0, 0)
		err = dc.LoadFontFace(lm.fontPath, 20)
		if err != nil {
			return nil, errors.WithMessage(err, "could not load font")
		}

		bcWidth, bcHeight := 50, 50

		br, err := qr.Encode(fmt.Sprintf("%s-%d/%d", strings.ToLower(shipment.Code), ul.Sequence, ulCount), qr.H, qr.Auto)
		if err != nil {
			return nil, errors.WithMessage(err, "could not encode qr code")
		}

		bc, err := barcode.Scale(br, bcWidth, bcHeight)
		if err != nil {
			return nil, errors.WithMessage(err, "could not scale barcode")
		}

		dc.DrawImage(bc, 10, 10)
		dc.DrawStringAnchored(fmt.Sprintf("%s-%d/%d", strings.ToUpper(shipment.Code), ul.Sequence, ulCount), 70, 45, 0, 0)
		dc.SetLineWidth(1)
		dc.DrawLine(10, 70, PAPER_W-10, 70)
		dc.Stroke()
		dc.LoadFontFace(lm.fontPath, 16)
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
		dc.LoadFontFace(lm.fontPath, 12)
		dc.DrawString("Кросс Азия", 10, 380)
		dc.DrawString("Телефон 1: +7(812)309-73-97", 10, 395)
		dc.DrawString("Телефон 2: +86 136-9923-3755", 10, 410)
		dc.DrawString("Сайт: https://crossasia.ru", 10, 425)

		err = dc.EncodePNG(file)

		if err != nil {
			return nil, errors.WithMessage(err, "could not encode png file")
		}

		ulBarcode.File = file
		ulBarcode.FullPath = tmpFilePath

		barcodes = append(barcodes, ulBarcode)
	}

	return barcodes, nil
}

func (lm LabelManager) CreateShipmentEntriesLabel(shipment logistics.Shipment) (Label, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	err := pdf.AddTTFFont("noto-cjk", lm.fontPath)
	if err != nil {
		return Label{}, err
	}

	err = pdf.SetFont("noto-cjk", "", 20)
	if err != nil {
		return Label{}, err
	}

	pdf.SetY(50)

	basicInformation := []string{
		fmt.Sprintf("票号: %s", shipment.Code),
		fmt.Sprintf("客户号: %s", shipment.CustomerCode),
		fmt.Sprintf("入库数量: %d", len(shipment.Entries)),
		fmt.Sprintf("包装方式: %s", shipment.PackageMethodZh),
	}

	for _, l := range basicInformation {
		pdf.SetX(20)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 24)
	}

	pdf.SetY(pdf.GetY() + 15)
	pdf.SetLineWidth(0.5)
	pdf.RectFromUpperLeft(20, pdf.GetY(), gopdf.PageSizeA4.W-40, 1)
	pdf.SetY(pdf.GetY() + 15)

	pdf.SetY(pdf.GetY() + 10)
	err = pdf.SetFont("noto-cjk", "", 16)
	if err != nil {
		return Label{}, err
	}

	for i, entry := range shipment.Entries {
		pdf.RectFromUpperLeft(20, pdf.GetY(), 15, 15)
		pdf.SetX(40)
		text := fmt.Sprintf("%d. %s (%d)   %s  %s", i+1, entry.ID, entry.BoxQty, entry.Source, entry.TrackCode)
		err = pdf.Cell(nil, text)
		if err != nil {
			return Label{}, err
		}

		safeSetY(pdf, pdf.GetY()+20, 50)
	}

	tmpFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-ShipmentEntriesLabel.pdf", xid.New()))
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return Label{}, errors.WithMessage(err, "could not create tmp file")
	}
	defer file.Close()

	err = pdf.WritePdf(tmpFilePath)
	if err != nil {
		return Label{}, err
	}

	res := Label{
		File:     file,
		FullPath: tmpFilePath,
	}

	return res, nil
}

func (lm LabelManager) CreateShipmentPartnerInfoLabel(shipment logistics.Shipment) (Label, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	err := pdf.AddTTFFont("noto-cjk", lm.fontPath)
	if err != nil {
		return Label{}, err
	}

	err = pdf.SetFont("noto-cjk", "", 16)
	if err != nil {
		return Label{}, err
	}

	pdf.SetY(50)

	title := fmt.Sprintf("票号 %s 出货信息", shipment.Code)
	titleWidth, err := pdf.MeasureTextWidth(title)

	if err != nil {
		return Label{}, nil
	}

	pdf.SetX(gopdf.PageSizeA4.W/2 - (titleWidth / 2))
	err = pdf.Text(title)
	if err != nil {
		return Label{}, nil
	}

	basicInfoStartY := pdf.GetY() + 24
	pdf.SetY(basicInfoStartY)

	basicInformation := []string{
		fmt.Sprintf("出货日期: %s", time.Now().Format("2006.01.02")),
		fmt.Sprintf("合作方: %s", shipment.PartnerInfo.Code),
		fmt.Sprintf("运输方式: %s", i18n.TranslateDeliveryMethod(shipment.PartnerInfo.DeliveryMethod)),
		fmt.Sprintf("保险: %s", i18n.TranslateCargoValueZh(shipment.PartnerInfo.CargoValue)),
		fmt.Sprintf("是否报关: %s", i18n.TranslateBoolZh(shipment.NeedDeclare)),
	}

	for _, l := range basicInformation {
		pdf.SetX(20)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 20)
	}

	cargoInfo := []string{
		fmt.Sprintf("总重量: %v kg", shipment.Weight()),
		fmt.Sprintf("总体积: %v m3", shipment.Cubage()),
		fmt.Sprintf("密度: %v kg/m3", shipment.Density()),
		fmt.Sprintf("箱数: %d 箱", len(shipment.UnitLoads)),
	}

	pdf.SetY(basicInfoStartY)

	for _, l := range cargoInfo {
		pdf.SetX(gopdf.PageSizeA4.W / 2)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 20)
	}

	recipientInfo := []string{
		fmt.Sprintf("收货人：%s", shipment.PartnerInfo.Recipient.Name),
		fmt.Sprintf("电话：%s", shipment.PartnerInfo.Recipient.PhoneNumber),
		fmt.Sprintf("目的地：%s", shipment.PartnerInfo.Recipient.Destination),
	}

	pdf.SetY(pdf.GetY() + 30)

	for _, l := range recipientInfo {
		pdf.SetX(20)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 20)
	}

	pdf.SetX(20)
	pdf.SetY(pdf.GetY() + 15)
	productName := fmt.Sprintf("品名: %s", shipment.PartnerInfo.ProductName)
	splitProductNames, err := pdf.SplitText(productName, gopdf.PageSizeA4.W-40)
	if err != nil {
		return Label{}, err
	}
	pnWidth, err := pdf.MeasureTextWidth("品名: ")
	if err != nil {
		return Label{}, nil
	}

	for i, pn := range splitProductNames {
		if i > 0 {
			pdf.SetX(20 + pnWidth)
		} else {
			pdf.SetX(20)
		}

		err = pdf.Cell(nil, pn)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 20)
	}

	pdf.SetY(pdf.GetY() + 15)
	pdf.SetLineWidth(0.5)
	pdf.RectFromUpperLeft(20, pdf.GetY(), gopdf.PageSizeA4.W-40, 1)
	pdf.SetY(pdf.GetY() + 20)

	for i, ul := range shipment.UnitLoads {
		pdf.SetX(20)
		weight := fmt.Sprintf("%d. %v kg", i+1, ul.Weight)
		size := fmt.Sprintf("%d × %d × %d", ul.Length, ul.Width, ul.Height)
		cubage := fmt.Sprintf("%v m3", ul.Cubage())
		err = pdf.Cell(nil, weight)
		if err != nil {
			return Label{}, err
		}

		pdf.SetX(140)
		err = pdf.Cell(nil, size)
		if err != nil {
			return Label{}, err
		}

		pdf.SetX(280)
		err = pdf.Cell(nil, cubage)
		if err != nil {
			return Label{}, err
		}

		pdf.SetX(370)
		pns, err := pdf.SplitText(ul.ProductName, gopdf.PageSizeA4.W-20-pdf.GetX())
		if err != nil {
			return Label{}, err
		}
		for i, pn := range pns {
			pdf.SetX(370)
			err = pdf.Cell(nil, pn)
			if err != nil {
				return Label{}, err
			}
			if len(pns)-1 > i {
				safeSetY(pdf, pdf.GetY()+20, 50)
			}
		}

		safeSetY(pdf, pdf.GetY()+20, 50)
	}

	tmpFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-ShipmentPartnerInfoLabel.pdf", xid.New()))
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return Label{}, errors.WithMessage(err, "could not create tmp file")
	}
	defer file.Close()

	err = pdf.WritePdf(tmpFilePath)
	if err != nil {
		return Label{}, err
	}

	res := Label{
		File:     file,
		FullPath: tmpFilePath,
	}

	return res, nil
}

func safeSetY(pdf *gopdf.GoPdf, y float64, newPageY float64) {
	if pdf.GetY()+100 > gopdf.PageSizeA4.H {
		pdf.AddPage()
		pdf.SetY(newPageY)
	} else {
		pdf.SetY(y)
	}
}
