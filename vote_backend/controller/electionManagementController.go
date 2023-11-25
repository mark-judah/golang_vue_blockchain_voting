package controller

import (
	"crypto/sha256"
	"encoding/hex"
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

func NewCounty(context *gin.Context) {
	var newCounty models.County
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
	} else {
		if err := context.BindJSON(&newCounty); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
	}
	token := context.GetHeader("Authorization")
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	email := GetCurrentUser(cleanToken)
	fmt.Println("token" + cleanToken)
	if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
		log.Fatalln(err)
	}
	if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "election-admin" {
		context.IndentedJSON(401, "You are not authorized to perform this action")
		return
	} else {
		result := database.Create(&newCounty)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newCounty)
		}
	}
}

func FetchCounties(context *gin.Context) {
	allCounties := []models.County{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allCounties).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allCounties))

		context.IndentedJSON(http.StatusOK, allCounties)
	}
}

func NewConstituency(context *gin.Context) {
	var newConstituency models.Constituency
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
	} else {
		if err := context.BindJSON(&newConstituency); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
	}
	token := context.GetHeader("Authorization")
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	email := GetCurrentUser(cleanToken)
	fmt.Println("token" + cleanToken)
	if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
		log.Fatalln(err)
	}
	if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "election-admin" {
		context.IndentedJSON(401, "You are not authorized to perform this action")
		return
	} else {
		result := database.Create(&newConstituency)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newConstituency)
		}
	}
}

func FetchConstituencies(context *gin.Context) {
	allConstituencies := []models.Constituency{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allConstituencies).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allConstituencies))

		context.IndentedJSON(http.StatusOK, allConstituencies)
	}
}

func NewWard(context *gin.Context) {
	var newWard models.Ward
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
	} else {
		if err := context.BindJSON(&newWard); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
	}
	token := context.GetHeader("Authorization")
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	email := GetCurrentUser(cleanToken)
	fmt.Println("token" + cleanToken)
	if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
		log.Fatalln(err)
	}
	if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "election-admin" {
		context.IndentedJSON(401, "You are not authorized to perform this action")
		return
	} else {
		result := database.Create(&newWard)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newWard)
		}
	}
}

func FetchWards(context *gin.Context) {
	allWards := []models.Ward{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allWards).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allWards))

		context.IndentedJSON(http.StatusOK, allWards)
	}
}

func NewPollingStation(context *gin.Context) {
	var newPollingStation models.PollingStation
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
	} else {
		if err := context.BindJSON(&newPollingStation); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
	}
	token := context.GetHeader("Authorization")
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	email := GetCurrentUser(cleanToken)
	fmt.Println("token" + cleanToken)
	if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
		log.Fatalln(err)
	}
	if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "election-admin" {
		context.IndentedJSON(401, "You are not authorized to perform this action")
		return
	} else {
		result := database.Create(&newPollingStation)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newPollingStation)
		}
	}
}

func FetchPollingStations(context *gin.Context) {
	allPollingStations := []models.PollingStation{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allPollingStations).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allPollingStations))

		context.IndentedJSON(http.StatusOK, allPollingStations)
	}
}

func NewDesktopClient(context *gin.Context) {
	var newDesktopClient models.DesktopClient
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
	} else {
		if err := context.BindJSON(&newDesktopClient); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
	}
	token := context.GetHeader("Authorization")
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	email := GetCurrentUser(cleanToken)
	fmt.Println("token" + cleanToken)
	if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
		log.Fatalln(err)
	}
	if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "election-admin" {
		context.IndentedJSON(401, "You are not authorized to perform this action")
		return
	} else {
		result := database.Create(&newDesktopClient)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newDesktopClient)
		}
	}
}

func FetchDesktopClients(context *gin.Context) {
	allDesktopClients := []models.DesktopClient{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allDesktopClients).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allDesktopClients))

		context.IndentedJSON(http.StatusOK, allDesktopClients)
	}
}

func NewCandidate(context *gin.Context) {
	var newCandidate models.Candidate
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
	} else {
		if err := context.BindJSON(&newCandidate); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
	}
	token := context.GetHeader("Authorization")
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	email := GetCurrentUser(cleanToken)
	fmt.Println("token" + cleanToken)
	if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
		log.Fatalln(err)
	}
	if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "election-admin" {
		context.IndentedJSON(401, "You are not authorized to perform this action")
		return
	} else {
		result := database.Create(&newCandidate)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newCandidate)
		}
	}
}

func FetchCandidates(context *gin.Context) {
	allCandidates := []models.Candidate{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allCandidates).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allCandidates))

		context.IndentedJSON(http.StatusOK, allCandidates)
	}
}

func NewVoter(context *gin.Context) {
	var newVoter models.Voter
	var currentUserDetails models.Users

	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)
	} else {
		if err := context.BindJSON(&newVoter); err != nil {
			context.IndentedJSON(http.StatusBadRequest, err.Error)
			context.Abort()
			return
		}
	}
	token := context.GetHeader("Authorization")
	cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
	email := GetCurrentUser(cleanToken)
	fmt.Println("token" + cleanToken)
	if err := database.Where("email=?", email).First(&currentUserDetails).Error; err != nil {
		log.Fatalln(err)
	}
	if currentUserDetails.Role != "superuser" && currentUserDetails.Role != "election-officer" {
		context.IndentedJSON(401, "You are not authorized to perform this action")
		return
	} else {
		newVoterDetailsBytes, err3 := json.Marshal(newVoter)
		if err3 != nil {
			panic(err3)
		}
		sum := sha256.Sum256([]byte(newVoterDetailsBytes))
		newVoterHash := hex.EncodeToString(sum[:])

		newVoter.VoterDetailsHash = newVoterHash
		result := database.Create(&newVoter)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newVoter)
		}
	}
}

func FetchVoters(context *gin.Context) {
	allVoters := []models.Voter{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Find(&allVoters).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(allVoters))

		context.IndentedJSON(http.StatusOK, allVoters)
	}
}

func FetchTransactionPool(context *gin.Context) {
}
