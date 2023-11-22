package models

import "time"

type Tally struct {
	CandidateId string `gorm:"primaryKey"`
	// BlockHeight string
	Total     int `gorm:"primaryKey"`
	Timestamp time.Time
}
