//go:generate gomodifytags -file $GOFILE -struct PartnerInfo -add-tags json -w
//go:generate gomodifytags -file $GOFILE -struct Recipient -add-tags json -w

package logistics

type PartnerInfo struct {
	Code           string         `json:"code"`
	DeliveryMethod DeliveryMethod `json:"transport_method"`
	ProductName    string         `json:"product_name"`
	Recipient      Recipient      `json:"recipient"`
	CargoValue     float64        `json:"cargo_value"`
}

type Recipient struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Destination string `json:"destination"`
}

type DeliveryMethod string

const (
	DMParcelExpress   DeliveryMethod = "parcel_express"
	DMAirExpress      DeliveryMethod = "air_express"
	DMAirEconomy      DeliveryMethod = "air_economy"
	DMLandRail        DeliveryMethod = "land_rail"
	DMLandRoadExpress DeliveryMethod = "land_road_express"
	DMLandRoadEconomy DeliveryMethod = "land_road_economy"
	DMLandContainer   DeliveryMethod = "land_container"
	DMWater           DeliveryMethod = "water"
	DMLandRoadCommon  DeliveryMethod = "land_road_common"
)
