package models

import (
	"gorm.io/gorm"
)

type Candidate struct {
	gorm.Model
	Name             string `json:"name" gorm:"unique"`
	Position         string
	CountyID         int `gorm:"column:county_id"`
	ConstituencyID   int `gorm:"column:constituency_id"`
	WardID           int `gorm:"column:ward_id"`
	PollingStationID int `gorm:"column:polling_station_id"`
}
