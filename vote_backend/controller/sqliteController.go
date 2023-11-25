package controller

import (
	"fmt"
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
	database.AutoMigrate(&models.Users{})
	database.AutoMigrate(&models.County{})
	database.AutoMigrate(&models.Constituency{})
	database.AutoMigrate(&models.Ward{})
	database.AutoMigrate(&models.PollingStation{})
	database.AutoMigrate(&models.DesktopClient{})
	database.AutoMigrate(&models.Candidate{})
	database.AutoMigrate(&models.Voter{})

	//add superuser with role as the primary key
	newUser := models.Users{
		Name:     "admin",
		Email:    "admin@superuser.com",
		Contact:  "072222222",
		Role:     "superuser",
		Password: "123456",
	}

	if err := newUser.HashPassword(newUser.Password); err != nil {
		panic(err)
	}

	result := database.Create(&newUser) // pass pointer of data to Create
	if result.Error != nil {
		fmt.Println(result.Error)
	}

}
