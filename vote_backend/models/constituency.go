package models

import (
	"gorm.io/gorm"
)

type Constituency struct {
	gorm.Model
	Name     string `json:"name" gorm:"unique"`
	CountyID int    `gorm:"column:county_id"`
}
