package models

import "time"

type Block struct {
	Index             int `gorm:"primaryKey"`
	Version           int
	PreviousBlockHash string
	CreatedBy         string
	CreatedAt         time.Time
	Data              string
}
