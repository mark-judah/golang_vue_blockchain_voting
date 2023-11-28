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
	"gorm.io/gorm/clause"
)

var NodeStats []models.NodeStats

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
		result := database.Omit(clause.Associations).Create(&newCounty)
		if result.Error != nil {
			context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
		} else {
			context.IndentedJSON(http.StatusCreated, newCounty)
			newAdminDashLog := models.AdminDashLog{Type: "County", Payload: newCounty}
			PersistAdminDashLog(newAdminDashLog)

			mqttMessage := models.Message{Type: "new_county", Payload: newCounty}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
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
			newAdminDashLog := models.AdminDashLog{Type: "Constituency", Payload: newConstituency}
			PersistAdminDashLog(newAdminDashLog)
			mqttMessage := models.Message{Type: "new_constituency", Payload: newConstituency}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
		}
	}
}

func FetchConstituencies(context *gin.Context) {
	county := []models.County{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Preload("Constituency").Find(&county).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(county))

		context.IndentedJSON(http.StatusOK, county)
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
			newAdminDashLog := models.AdminDashLog{Type: "Ward", Payload: newWard}
			PersistAdminDashLog(newAdminDashLog)
			mqttMessage := models.Message{Type: "new_ward", Payload: newWard}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
		}
	}
}

func FetchWards(context *gin.Context) {
	county := []models.County{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Preload("Constituency.Ward").Find(&county).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(county))

		context.IndentedJSON(http.StatusOK, county)
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
			newAdminDashLog := models.AdminDashLog{Type: "PollingStation", Payload: newPollingStation}
			PersistAdminDashLog(newAdminDashLog)
			mqttMessage := models.Message{Type: "new_polling_station", Payload: newPollingStation}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
		}
	}
}

func FetchPollingStations(context *gin.Context) {
	county := []models.County{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Preload("Constituency.Ward.PollingStation").Find(&county).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(county))

		context.IndentedJSON(http.StatusOK, county)
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
			newAdminDashLog := models.AdminDashLog{Type: "DesktopClient", Payload: newDesktopClient}
			PersistAdminDashLog(newAdminDashLog)
			mqttMessage := models.Message{Type: "new_desktop_client", Payload: newDesktopClient}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
		}
	}
}

func FetchDesktopClients(context *gin.Context) {
	county := []models.County{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Preload("Constituency.Ward.PollingStation.DesktopClient").Find(&county).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(county))

		context.IndentedJSON(http.StatusOK, county)
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
			newAdminDashLog := models.AdminDashLog{Type: "Candidate", Payload: newCandidate}
			PersistAdminDashLog(newAdminDashLog)
			mqttMessage := models.Message{Type: "new_candidate", Payload: newCandidate}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
		}
	}
}

func FetchCandidates(context *gin.Context) {
	county := []models.County{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Preload("Constituency.Ward.PollingStation.Candidate").Find(&county).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(county))

		context.IndentedJSON(http.StatusOK, county)
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
			newAdminDashLog := models.AdminDashLog{Type: "Voter", Payload: newVoter}
			PersistAdminDashLog(newAdminDashLog)
			mqttMessage := models.Message{Type: "new_voter", Payload: newVoter}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := Client[0].Publish("adminTransaction/1", 0, false, data)
			token.Wait()
		}
	}
}

func FetchVoters(context *gin.Context) {
	county := []models.County{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, err)

	} else {
		if err := database.Preload("Constituency.Ward.PollingStation.Voter").Find(&county).Error; err != nil {
			log.Fatalln(err)
		}
		log.Printf("%d rows found.", len(county))

		context.IndentedJSON(http.StatusOK, county)
	}
}

func FetchTransactionPool(context *gin.Context) {
}

func FetchConnectedNodes(context *gin.Context) {
	token := Client[0].Publish("nodeStatsRequest/1", 0, false, "get node stats")
	token.Wait()

	for {
		if len(NodeStats) < getTotalConnectedNodes() {
			fmt.Println("fetching node stats....")
		} else {
			context.IndentedJSON(http.StatusOK, NodeStats)
			NodeStats = nil
			break
		}
	}
}
