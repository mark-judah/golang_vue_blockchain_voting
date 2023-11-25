package models

import (
	"gorm.io/gorm"
)

type Ward struct {
	gorm.Model
	Name           string `json:"name" gorm:"unique"`
	CountyID       int    `gorm:"column:county_id"`
	ConstituencyID int    `gorm:"column:constituency_id"`
}
