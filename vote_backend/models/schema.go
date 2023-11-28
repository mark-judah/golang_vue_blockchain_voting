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

type Candidate struct {
	gorm.Model
	Name             string `json:"name" gorm:"unique"`
	Position         string
	CountyID         int `gorm:"column:county_id"`
	ConstituencyID   int `gorm:"column:constituency_id"`
	WardID           int `gorm:"column:ward_id"`
	PollingStationID int `gorm:"column:polling_station_id"`
}

type PollingStation struct {
	gorm.Model
	Name           string `json:"name" gorm:"unique"`
	CountyID       int    `gorm:"column:county_id"`
	ConstituencyID int    `gorm:"column:constituency_id"`
	WardID         int    `gorm:"column:ward_id"`
	Candidate      []Candidate
	Voter          []Voter
	DesktopClient  []DesktopClient
}

type Ward struct {
	gorm.Model
	Name           string `json:"name" gorm:"unique"`
	CountyID       int    `gorm:"column:county_id"`
	ConstituencyID int    `gorm:"column:constituency_id"`
	PollingStation []PollingStation
}

type Constituency struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique"`
	CountyID int    `gorm:"column:county_id"`
	Ward     []Ward
}

type County struct {
	gorm.Model
	Name         string `json:"name" gorm:"unique"`
	Constituency []Constituency
}
