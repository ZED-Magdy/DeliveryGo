package Dtos

type PriceListDto struct {
	KmCost          float64 `json:"km_cost"`
	CancellationCost float64 `json:"cancellation_cost"`
}