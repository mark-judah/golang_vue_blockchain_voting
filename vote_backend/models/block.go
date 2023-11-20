package models

import "time"

type Block struct {
	Index             int `gorm:"primaryKey;auto_increment;not_null"`
	Version           int
	BlockHash         string
	PreviousBlockHash string
	CreatedBy         string
	CreatedAt         time.Time
	Data              string
}
