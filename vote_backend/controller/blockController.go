package controller

import (
	"encoding/json"
	"fmt"
	"vote_backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateBlock() {
	fmt.Println("creating block")

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	//fetch oldest five rows
	var transactions []models.Transaction
	database.Limit(5).Order("created_at asc").Find(&transactions)

	jsonData, err := json.Marshal(transactions)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonData))

	//create block

	//insert into db
}
