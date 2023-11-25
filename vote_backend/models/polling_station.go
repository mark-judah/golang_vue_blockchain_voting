package models

import (
	"gorm.io/gorm"
)

type PollingStation struct {
	gorm.Model
	Name           string `json:"name" gorm:"unique"`
	CountyID       int    `gorm:"column:county_id"`
	ConstituencyID int    `gorm:"column:constituency_id"`
	WardID         int    `gorm:"column:ward_id"`
}
