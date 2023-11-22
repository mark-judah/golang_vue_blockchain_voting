package controller

import (
	"vote_backend/models"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSqlite() {
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	database.AutoMigrate(&models.Transaction{})
	database.AutoMigrate(&models.Block{})
	database.AutoMigrate(&models.Tally{})

}
