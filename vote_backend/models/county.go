package models

import (
	"gorm.io/gorm"
)

type County struct {
	gorm.Model
	Name string `json:"name" gorm:"unique"`
}
