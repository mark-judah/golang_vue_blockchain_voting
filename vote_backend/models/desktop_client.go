package models

import (
	"gorm.io/gorm"
)

type DesktopClient struct {
	gorm.Model
	Name             string `json:"name" gorm:"unique"`
	CountyID         int    `gorm:"column:county_id"`
	ConstituencyID   int    `gorm:"column:constituency_id"`
	WardID           int    `gorm:"column:ward_id"`
	PollingStationID int    `gorm:"column:polling_station_id"`
}
