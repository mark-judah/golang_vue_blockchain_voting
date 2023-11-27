package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"vote_backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateUser(context *gin.Context) {
	var newUser models.Users
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := context.BindJSON(&newUser); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
		if err := newUser.HashPassword(newUser.Password); err != nil {
			context.IndentedJSON(http.StatusInternalServerError, err.Error)
			context.Abort()
			return
		}
		token := context.GetHeader("Authorization")
		cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
		email := GetCurrentUser(cleanToken)
		fmt.Println("token" + cleanToken)
		if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
			log.Fatalln(err)
		}
		switch {
		case newUser.Role == "superuser":
			context.IndentedJSON(409, "User already exists")
			return

		case newUser.Role == "hr-admin":
			if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "hr-admin" {
				context.IndentedJSON(401, "You are not authorized to perform this action")
				return
			}
		case newUser.Role == "election-admin":
			if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "hr-admin" {
				context.IndentedJSON(401, "You are not authorized to perform this action")
				return
			}
		case newUser.Role == "election-officer":
			if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "hr-admin" {
				context.IndentedJSON(401, "You are not authorized to perform this action")
				return
			}
		case newUser.Role != "superuser" && newUser.Role != "hr-admin" && newUser.Role != "election-admin" && newUser.Role != "election-officer":
			context.IndentedJSON(400, "Invalid user role")
			return
		}

		result := database.Create(&newUser) // pass pointer of data to Create
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newUser)
			newAdminLog := models.AdminDashLog{Type: "User", Payload: newUser}
			PersistAdminDashLog(newAdminLog)
			mqttMessage := models.Message{Type: "new_user", Payload: newUser}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
		}

	}
}

func GetUsers(context *gin.Context) {
	allUsers := []models.Users{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allUsers).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allUsers))

		context.IndentedJSON(http.StatusOK, allUsers)
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(context *gin.Context) {
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		var loginRequest LoginRequest
		var user models.Users

		if err := context.BindJSON(&loginRequest); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
		// check if email exists and password is correct
		result := database.Where("email=?", loginRequest.Email).First(&user)
		if result.Error != nil {
			context.IndentedJSON(http.StatusInternalServerError, result.Error)
			context.Abort()
			return
		}
		credentialError := user.CheckPassword(loginRequest.Password)
		if credentialError != nil {
			context.IndentedJSON(http.StatusBadRequest, "Invalid Credentials")
			context.Abort()
			return
		}
		token, err := GenerateJWT(user.Email, user.Name)
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, err)
			context.Abort()
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"token":    token,
			"email":    user.Email,
			"username": user.Name,
		})
	}
}

func CurrentUser(context *gin.Context) {

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		var user models.Users
		token := context.GetHeader("Authorization")
		cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
		email := GetCurrentUser(cleanToken)

		if err := database.Where("email=?", email).First(&user).Error; err != nil {
			log.Fatalln(err)
		}

		context.IndentedJSON(http.StatusOK, user)

	}
}
