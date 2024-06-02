package Dtos

type InitInfoResponseDto struct {
	Region RegionDto `json:"region"`
	Geofence GeofenceDto `json:"geofence"`
	PriceList PriceListDto `json:"price_list"`
}