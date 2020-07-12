package i18n

import (
	"fmt"
	"github.com/amanbolat/ca-warehouse-client/logistics"
)

func TranslateBoolZh(v bool) string {
	if v {
		return "是"
	}

	return "否"
}

func TranslateCargoValueZh(v float64) string {
	if v == 0 {
		return "义务保险"
	}

	return fmt.Sprintf("%.0f 美元", v)
}

func TranslateDeliveryMethod(dm logistics.DeliveryMethod) string {
	switch dm {
	case logistics.DMAirEconomy:
		return "慢空"
	case logistics.DMAirExpress:
		return "快空"
	case logistics.DMLandContainer:
		return "陆运集装箱"
	case logistics.DMLandRail:
		return "铁路"
	case logistics.DMLandRoadCommon:
		return "普通汽运"
	case logistics.DMLandRoadEconomy:
		return "慢汽运"
	case logistics.DMLandRoadExpress:
		return "快汽运"
	case logistics.DMParcelExpress:
		return "速递快递"
	case logistics.DMWater:
	}

	return "无知"
}
