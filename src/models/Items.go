package models

import "gorm.io/gorm"

type Items struct {
	ID          uint   `gorm:"primary_key" json:"id"`
	NamaBarang  string `json:"nama_barang"`
	HargaBarang int    `json:"harga_barang"`
	UserId      uint   `json:"user_id"`
	gorm.Model

	Users Users `gorm:"foreignkey:UserId"`
}
