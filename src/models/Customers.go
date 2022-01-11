package models

import "gorm.io/gorm"

type Customer struct {
	ID          uint   `gorm:"primary_key"`
	UserID      uint   `gorm:"not null"`
	Nama        string `gorm:"not null"`
	PhoneNumber string `gorm:"not null"`
	gorm.Model

	Users Users `gorm:"foreignkey:UserID"`
}
