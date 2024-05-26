package Models

import (

	"gorm.io/gorm"
)

type PriceList struct {
	gorm.Model
	Name             string
	KmCost           float64
	CancellationCost float64
}
