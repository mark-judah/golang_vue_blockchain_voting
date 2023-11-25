package models

import "gorm.io/gorm"

type Voter struct {
	gorm.Model
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	VoterId          string `json:"voterId" gorm:"primaryKey"`
	PhoneNumber      string `json:"phoneNumber" gorm:"primaryKey"`
	VoterDetailsHash string `json:"voterDetailsHash"`
	CountyID         int    `gorm:"column:county_id"`
	ConstituencyID   int    `gorm:"column:constituency_id"`
	WardID           int    `gorm:"column:ward_id"`
	PollingStationID int    `gorm:"column:polling_station_id"`
}
