package models

import "gorm.io/gorm"

type Notifications struct {
	ID               uint   `gorm:"primary_key" json:"id"`
	TransactionId    uint   `json:"transaction_id"`
	NotificationTeks string `json:"notification_teks"`
	HasSent          bool   `json:"has_sent";default:"false"`
	PhoneNumber      string `json:"phone_number"`

	gorm.Model
	Transaksi Transaksi `gorm:"foreignkey:TransactionId"`
}
