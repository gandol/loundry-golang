package models

import "gorm.io/gorm"

type TransaksiItems struct {
	ID          uint   `gorm:"primary_key"`
	TransaksiID uint   `gorm:"column:transaksi_id"`
	IdItem      uint   `gorm:"column:id_item"`
	Qty         int    `gorm:"column:qty";default:0`
	Satuan      string `gorm:"column:satuan";default:""`
	Total       int    `gorm:"column:total";default:0`
	StatusItem  string `gorm:"column:status_item";default:""`
	gorm.Model

	Transaksi Transaksi `gorm:"foreignkey:TransaksiID"`
	Items     Items     `gorm:"foreignkey:IdItem"`
}
