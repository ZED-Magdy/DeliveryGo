package Models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string
	Email     *string `gorm:"unique"`
	Phone     string  `gorm:"unique"`
	Password  string
	Otp       *string
	CurrentRegionId uint       `gorm:"foreignKey:Id,constraint:OnDelete:SET NULL"`
	CurrentRegion   Region  `gorm:"foreignKey:CurrentRegionId"`
}