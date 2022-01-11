package models

import "gorm.io/gorm"

type Transaksi struct {
	ID         uint    `gorm:"primary_key"`
	UserId     uint    `gorm:"not null"`
	CustomerId uint    `gorm:"not null"`
	SubTotal   float64 `gorm:"default:0"`
	Diskon     float64 `gorm:"default:0"`
	Total      float64 `gorm:"default:0"`
	Status     string  `gorm:"not null"`

	gorm.Model

	Users    Users    `gorm:"foreignkey:UserId"`
	Customer Customer `gorm:"foreignkey:CustomerId"`
}
