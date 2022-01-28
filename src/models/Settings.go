package models

import "gorm.io/gorm"

type Settings struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	UserID      uint   `json:"user_id"`
	gorm.Model
	Users Users `gorm:"foreignkey:UserID"`
}
