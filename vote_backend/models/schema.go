package models

import (
	"time"

	"gorm.io/gorm"
)

type Block struct {
	Index             int `gorm:"primaryKey;auto_increment;not_null"`
	Version           int
	BlockHash         string
	PreviousBlockHash string
	CreatedBy         string
	CreatedAt         time.Time
	Data              string
}

type Transaction struct {
	Txid             string    `json:"txid" gorm:"primaryKey"`
	NodeId           string    `json:"nodeId"`
	CandidateId      string    `json:"candidateId"`
	CountyID         int       `json:"county_id" gorm:"column:county_id"`
	ConstituencyID   int       `json:"constituency_id" gorm:"column:constituency_id"`
	WardID           int       `json:"ward_id" gorm:"column:ward_id"`
	PollingStationID int       `json:"polling_station_id" gorm:"column:polling_station_id"`
	CreatedAt        time.Time `json:"timestamp"`
}

type DesktopClient struct {
	gorm.Model
	SerialNumber     string `json:"serial_number" gorm:"unique"`
	CountyID         int    `json:"county_id" gorm:"column:county_id"`
	ConstituencyID   int    `json:"constituency_id" gorm:"column:constituency_id"`
	WardID           int    `json:"ward_id" gorm:"column:ward_id"`
	PollingStationID int    `json:"polling_station_id" gorm:"column:polling_station_id"`
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
type Tally struct {
	gorm.Model
	CandidateID string `gorm:"primaryKey"`
	// BlockHeight string
	Total     int `gorm:"primaryKey"`
	Timestamp time.Time
}

type Candidate struct {
	gorm.Model
	Name             string `json:"name" gorm:"unique"`
	Position         string `json:"position"`
	Party            string `json:"party"`
	Slogan           string `json:"slogan"`
	Statement        string `json:"statement"`
	Photo            string `json:"photo"`
	CountyID         int    `gorm:"column:county_id"`
	ConstituencyID   int    `gorm:"column:constituency_id"`
	WardID           int    `gorm:"column:ward_id"`
	PollingStationID int    `gorm:"column:polling_station_id"`
	Tally            Tally
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
