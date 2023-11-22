package models

import "time"

type Tally struct {
	Index       int `gorm:"primaryKey;auto_increment;not_null"`
	CandidateId string
	// BlockHeight string
	Total     int
	Timestamp time.Time
}
