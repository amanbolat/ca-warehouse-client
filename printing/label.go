package printing

import (
	"bytes"
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/i18n"
	"github.com/amanbolat/ca-warehouse-client/logistics"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/signintech/gopdf"
	"image/jpeg"
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

func (lm LabelManager) CreateUnitLoadLabels(shipment logistics.Shipment) (*Label, error) {
	ulCount := len(shipment.UnitLoads)
	if ulCount < 1 {
		return nil, errors.New("there is now unit loads, nothing to print")
	}

	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{
		W: PAPER_W,
		H: PAPER_H,
	}})
	err := pdf.AddTTFFont("noto-cjk", lm.fontPath)
	err = pdf.SetFont("noto-cjk", "", 20)
	if err != nil {
		return nil, err
	}

	for _, ul := range shipment.UnitLoads {
		pdf.AddPage()
		// SPN00123
		pdf.SetY(50)
		err = pdf.SetFont("noto-cjk", "", 48)
		if err != nil {
			return nil, err
		}
		err = writeCenteredText(pdf, strings.ToUpper(shipment.Code))
		if err != nil {
			return nil, err
		}

		// 1/1
		err = pdf.SetFont("noto-cjk", "", 28)
		if err != nil {
			return nil, err
		}
		pdf.SetY(pdf.GetY() + 40)
		err = writeCenteredText(pdf, fmt.Sprintf("%d/%d", ul.Sequence, ulCount))
		if err != nil {
			return nil, err
		}

		err = pdf.SetFont("noto-cjk", "", 24)
		if err != nil {
			return nil, err
		}

		// Unit load weight and cubage
		err = pdf.SetFont("noto-cjk", "", 24)
		if err != nil {
			return nil, err
		}
		pdf.RectFromUpperLeft(5, pdf.GetY()+25, PAPER_W-10, 90)
		weight := fmt.Sprintf("%v kg", ul.Weight)
		cubage := fmt.Sprintf("%v m3", ul.Cubage())

		weightW, err := pdf.MeasureTextWidth(weight)
		if err != nil {
			return nil, err
		}
		pdf.SetX(PAPER_W/4 - weightW/2)
		pdf.SetY(pdf.GetY() + 35)
		err = pdf.Cell(nil, weight)
		if err != nil {
			return nil, err
		}

		cubageW, err := pdf.MeasureTextWidth(cubage)
		if err != nil {
			return nil, err
		}
		pdf.SetX(PAPER_W - PAPER_W/4 - cubageW/2)
		pdf.SetY(pdf.GetY())
		err = pdf.Cell(nil, cubage)
		if err != nil {
			return nil, err
		}

		pdf.SetY(pdf.GetY() + 65)
		err = pdf.SetFont("noto-cjk", "", 18)
		if err != nil {
			return nil, err
		}
		size := fmt.Sprintf("%d cm × %d cm × %d cm", ul.Length, ul.Width, ul.Height)
		err = writeCenteredText(pdf, size)
		if err != nil {
			return nil, err
		}

		// Shipment weight and cubage
		err = pdf.SetFont("noto-cjk", "", 12)
		if err != nil {
			return nil, err
		}
		pdf.SetY(pdf.GetY() + 40)

		shipmentInfo := []string{
			"Габариты всего груза/货物整体规格:",
			fmt.Sprintf("ВЕС/重量: %v kg", shipment.Weight()),
			fmt.Sprintf("ОБЪЕМ/体积: %v m3", shipment.Cubage()),
		}

		for _, str := range shipmentInfo {
			pdf.SetX(5)
			pdf.SetY(pdf.GetY() + 18)
			err = pdf.Cell(nil, str)
			if err != nil {
				return nil, err
			}
		}

		// Attention
		err = pdf.SetFont("noto-cjk", "", 10)
		if err != nil {
			return nil, err
		}
		attentionText := []string{
			"ВНИМАНИЕ! При получении груза",
			"обязательно проверьте целостность",
			"внешней упаковки и характеристики",
			"груза!!!",
		}
		pdf.SetY(PAPER_H - 80)

		for _, str := range attentionText {
			pdf.SetX(5)
			pdf.SetY(pdf.GetY() + 14)
			err = pdf.Cell(nil, str)
			if err != nil {
				return nil, err
			}
		}

		br, err := qr.Encode(fmt.Sprintf("%s-%d/%d", strings.ToLower(shipment.Code), ul.Sequence, ulCount), qr.H, qr.Auto)
		if err != nil {
			return nil, errors.WithMessage(err, "could not encode qr code")
		}

		bc, err := barcode.Scale(br, 100, 100)
		if err != nil {
			return nil, errors.WithMessage(err, "could not scale barcode")
		}

		buf := bytes.NewBuffer([]byte{})
		err = jpeg.Encode(buf, bc, &jpeg.Options{Quality: 100})
		if err != nil {
			return nil, errors.WithStack(err)
		}

		img, err := gopdf.ImageHolderByReader(buf)
		if err != nil {
			return nil, err
		}
		err = pdf.ImageByHolder(img, PAPER_W-70, PAPER_H-70, nil)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	tmpFilePath := path.Join(os.TempDir(), fmt.Sprintf("%s-unitload-labels.pdf", xid.New()))
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return nil, errors.WithMessage(err, "could not create tmp file")
	}
	defer file.Close()

	err = pdf.WritePdf(tmpFilePath)
	if err != nil {
		return nil, err
	}

	label := &Label{
		File:     file,
		FullPath: tmpFilePath,
	}

	return label, nil
}

func (lm LabelManager) CreateShipmentPreparationLabels(shipment logistics.Shipment) (Label, error) {
	pdf := &gopdf.GoPdf{}
	pageSize := gopdf.Rect{
		W: PAPER_W,
		H: PAPER_H,
	}
	pdf.Start(gopdf.Config{PageSize: pageSize})
	pdf.AddPage()
	err := pdf.AddTTFFont("noto-cjk", lm.fontPath)
	if err != nil {
		return Label{}, err
	}

	// Basic information
	err = pdf.SetFont("noto-cjk", "", 16)
	if err != nil {
		return Label{}, err
	}
	pdf.SetY(15)

	basicInformation := []string{
		fmt.Sprintf("票号: %s", shipment.Code),
		fmt.Sprintf("客户号: %s", shipment.CustomerCode),
		fmt.Sprintf("入库数量: %d", len(shipment.Entries)),
		fmt.Sprintf("包装方式: %s", shipment.PackageMethodZh),
	}

	for _, l := range basicInformation {
		pdf.SetX(5)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 20)
	}

	// Line
	pdf.SetY(pdf.GetY() + 10)
	pdf.RectFromUpperLeft(5, pdf.GetY(), PAPER_W-10, 1)
	pdf.SetY(pdf.GetY() + 10)

	// Note list
	pdf.SetY(pdf.GetY() + 10)
	err = pdf.SetFont("noto-cjk", "", 16)
	if err != nil {
		return Label{}, err
	}
	err = writeCenteredText(pdf, "备注")
	if err != nil {
		return Label{}, err
	}

	err = pdf.SetFont("noto-cjk", "", 12)
	if err != nil {
		return Label{}, err
	}
	pdf.SetY(pdf.GetY() + 10)

	for i, note := range shipment.Notes {
		pdf.SetX(5)
		err = pdf.Cell(nil, fmt.Sprintf("%d: %s", i+1, note.Content))
		if err != nil {
			return Label{}, err
		}
		safeSetY(pdf, pdf.GetY()+16, 5)
	}

	// Line
	pdf.SetY(pdf.GetY() + 5)
	pdf.RectFromUpperLeft(5, pdf.GetY(), PAPER_W-10, 1)
	pdf.SetY(pdf.GetY() + 10)

	// Entry list
	pdf.SetY(pdf.GetY() + 10)
	err = pdf.SetFont("noto-cjk", "", 12)
	if err != nil {
		return Label{}, err
	}

	for i, entry := range shipment.Entries {
		// pdf.RectFromUpperLeft(5, pdf.GetY(), 13, 13)
		pdf.SetX(5)
		text := fmt.Sprintf("%d. %s (%d) %s %s", i+1, entry.ID, entry.BoxQty, entry.Source, entry.TrackCode)
		err = pdf.Cell(nil, text)
		if err != nil {
			return Label{}, err
		}

		safeSetY(pdf, pdf.GetY()+16, 5)
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
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{
		W: PAPER_W,
		H: PAPER_H,
	}})
	pdf.AddPage()
	err := pdf.AddTTFFont("noto-cjk", lm.fontPath)
	if err != nil {
		return Label{}, err
	}

	// Ttile
	err = pdf.SetFont("noto-cjk", "", 14)
	if err != nil {
		return Label{}, err
	}

	pdf.SetY(15)

	title := fmt.Sprintf("票号 %s 出货信息", shipment.Code)
	err = writeCenteredText(pdf, title)
	if err != nil {
		return Label{}, nil
	}

	// Basic Info
	err = pdf.SetFont("noto-cjk", "", 12)
	if err != nil {
		return Label{}, err
	}
	basicInfoStartY := pdf.GetY() + 15
	pdf.SetY(basicInfoStartY)

	basicInformation := []string{
		fmt.Sprintf("出货日期: %s", time.Now().Format("2006.01.02")),
		fmt.Sprintf("合作方: %s", shipment.PartnerInfo.Code),
		fmt.Sprintf("运输方式: %s", i18n.TranslateDeliveryMethod(shipment.PartnerInfo.DeliveryMethod)),
		fmt.Sprintf("保险: %s", i18n.TranslateCargoValueZh(shipment.PartnerInfo.CargoValue)),
		fmt.Sprintf("是否报关: %s", i18n.TranslateBoolZh(shipment.NeedDeclare)),
	}

	for _, l := range basicInformation {
		pdf.SetX(5)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 16)
	}

	cargoInfo := []string{
		fmt.Sprintf("总重量: %v kg", shipment.Weight()),
		fmt.Sprintf("总体积: %v m3", shipment.Cubage()),
		fmt.Sprintf("密度: %v kg/m3", shipment.Density()),
		fmt.Sprintf("箱数: %d 箱", len(shipment.UnitLoads)),
	}

	pdf.SetY(basicInfoStartY)

	for _, l := range cargoInfo {
		pdf.SetX(PAPER_W / 2)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 16)
	}

	// Recipient
	recipientInfo := []string{
		fmt.Sprintf("收货人：%s", shipment.PartnerInfo.Recipient.Name),
		fmt.Sprintf("电话：%s", shipment.PartnerInfo.Recipient.PhoneNumber),
		fmt.Sprintf("目的地：%s", shipment.PartnerInfo.Recipient.Destination),
	}

	pdf.SetY(pdf.GetY() + 30)

	for _, l := range recipientInfo {
		pdf.SetX(5)
		err = pdf.Cell(nil, l)
		if err != nil {
			return Label{}, err
		}
		pdf.SetY(pdf.GetY() + 16)
	}

	// Product name
	pdf.SetX(5)
	pdf.SetY(pdf.GetY() + 15)
	pnWidth, err := pdf.MeasureTextWidth("品名: ")
	if err != nil {
		return Label{}, err
	}
	productName := fmt.Sprintf("品名: %s", shipment.PartnerInfo.ProductName)
	splitProductNames := safeSplitText(pdf, productName, PAPER_W-10-pnWidth)

	for i, pn := range splitProductNames {
		if i > 0 {
			pdf.SetX(5 + pnWidth)
		} else {
			pdf.SetX(5)
		}

		err = pdf.Cell(nil, pn)
		if err != nil {
			return Label{}, errors.WithStack(err)
		}
		pdf.SetY(pdf.GetY() + 16)
	}

	// Line
	pdf.SetY(pdf.GetY() + 10)
	pdf.RectFromUpperLeft(5, pdf.GetY(), PAPER_W-10, 1)
	pdf.SetY(pdf.GetY() + 10)

	// Unit loads
	err = pdf.SetFont("noto-cjk", "", 12)
	if err != nil {
		return Label{}, err
	}

	for i, ul := range shipment.UnitLoads {
		pdf.SetX(5)
		weight := fmt.Sprintf("%d. %v kg", i+1, ul.Weight)
		size := fmt.Sprintf("%d×%d×%d", ul.Length, ul.Width, ul.Height)
		cubage := fmt.Sprintf("%v m3", ul.Cubage())
		err = pdf.Cell(nil, weight)
		if err != nil {
			return Label{}, err
		}

		pdf.SetX(120)
		err = pdf.Cell(nil, size)
		if err != nil {
			return Label{}, err
		}

		pdf.SetX(220)
		err = pdf.Cell(nil, cubage)
		if err != nil {
			return Label{}, err
		}

		// Product names
		sequenceW, err := pdf.MeasureTextWidth(fmt.Sprintf("%d. ", i+1))
		safeSetY(pdf, pdf.GetY()+16, 15)
		pns := safeSplitText(pdf, ul.ProductName, gopdf.PageSizeA4.W-20-pdf.GetX())
		for i, pn := range pns {
			pdf.SetX(5 + sequenceW)
			err = pdf.Cell(nil, pn)
			if err != nil {
				return Label{}, err
			}
			if len(pns)-1 > i {
				safeSetY(pdf, pdf.GetY()+16, 15)
			}
		}

		safeSetY(pdf, pdf.GetY()+16, 15)
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
	if pdf.GetY()+40 > PAPER_H {
		pdf.AddPage()
		pdf.SetY(newPageY)
	} else {
		pdf.SetY(y)
	}
}

func writeCenteredText(pdf *gopdf.GoPdf, text string) error {
	txtWidth, err := pdf.MeasureTextWidth(text)

	if err != nil {
		return err
	}

	pdf.SetX(PAPER_W/2 - (txtWidth / 2))
	err = pdf.Text(text)
	if err != nil {
		return err
	}

	return nil
}

func safeSplitText(pdf *gopdf.GoPdf, txt string, width float64) []string {
	if strings.TrimSpace(txt) == "" {
		return []string{}
	}
	strArr, _ := pdf.SplitText(txt, width)
	return strArr
}
