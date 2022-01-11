package models

import "gorm.io/gorm"

type Users struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"type:varchar(100)"`
	Nama     string `gorm:"type:varchar(100)"`
	gorm.Model
}
