package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

		type data struct {
			County       string
			Constituency string
		}
		var constituencies []data

		for _, a := range county {
			for _, b := range a.Constituency {
				constituencies = append(constituencies, data{County: a.Name, Constituency: b.Name})
			}
		}
		context.IndentedJSON(http.StatusOK, constituencies)
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
		type data struct {
			County       string
			Constituency string
			Ward         string
		}
		var wards []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					wards = append(wards, data{County: a.Name, Constituency: b.Name, Ward: c.Name})
				}
			}
		}
		context.IndentedJSON(http.StatusOK, wards)
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

		type data struct {
			County         string
			Constituency   string
			Ward           string
			PollingStation string
		}
		var polling_stations []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					for _, d := range c.PollingStation {
						polling_stations = append(polling_stations, data{County: a.Name, Constituency: b.Name, Ward: c.Name, PollingStation: d.Name})
					}
				}
			}
		}
		context.IndentedJSON(http.StatusOK, polling_stations)
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

		type data struct {
			County         string
			Constituency   string
			Ward           string
			PollingStation string
			Name           string
			SerialNumber   string
			MacAddress     string
		}
		var desktop_clients []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					for _, d := range c.PollingStation {
						for _, e := range d.DesktopClient {
							desktop_clients = append(desktop_clients, data{County: a.Name, Constituency: b.Name, Ward: c.Name, PollingStation: d.Name, Name: e.Name, SerialNumber: e.SerialNumber, MacAddress: e.MacAddress})
						}
					}
				}
			}
		}
		context.IndentedJSON(http.StatusOK, desktop_clients)
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

		type data struct {
			County         string
			Constituency   string
			Ward           string
			PollingStation string
			Candidate      string
			Position       string
			Party          string
		}
		var candidates []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					for _, d := range c.PollingStation {
						for _, e := range d.Candidate {
							candidates = append(candidates, data{County: a.Name, Constituency: b.Name, Ward: c.Name, PollingStation: d.Name, Candidate: e.Name, Position: e.Position, Party: e.Party})
						}
					}
				}
			}
		}
		context.IndentedJSON(http.StatusOK, candidates)
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

		type data struct {
			County         string
			Constituency   string
			Ward           string
			PollingStation string
			FirstName      string
			LastName       string
			IDNumber       string
			PhoneNumber    string
		}
		var voters []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					for _, d := range c.PollingStation {
						for _, e := range d.Voter {
							voters = append(voters, data{County: a.Name, Constituency: b.Name, Ward: c.Name, PollingStation: d.Name, FirstName: e.FirstName, LastName: e.LastName, IDNumber: e.VoterId, PhoneNumber: e.PhoneNumber})
						}
					}
				}
			}
		}
		context.IndentedJSON(http.StatusOK, voters)
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

func FetchQuickStats(context *gin.Context) {
	var totalVoters = 0
	var totalPollingStations = 0
	var totalVotes = 0
	var totalPresidentialCandidates = 0
	var totalDesktopClients = 0
	var totalOnlineClients = 0
	var transactionPoolSize = 0
	var totalProcessedTransactions = 0

	//get total voters
	voters := []models.Voter{}
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	if err := database.Find(&voters).Error; err != nil {
		log.Fatalln(err)
	}
	totalVoters = len(voters)

	//get total polling stations
	pollingStations := []models.PollingStation{}
	if err := database.Find(&pollingStations).Error; err != nil {
		log.Fatalln(err)
	}
	totalPollingStations = len(pollingStations)

	//get total votes
	blockChainFile, err := os.ReadFile("chain.json")
	if err != nil {
		totalVotes = 0
	} else {
		var nodesBlocks []models.Block

		err2 := json.Unmarshal(blockChainFile, &nodesBlocks)
		if err2 != nil {
			panic(err2)
		}
		var totalTransactions []models.Transaction

		for _, x := range nodesBlocks {
			data := x.Data
			var transaction []models.Transaction
			err2 := json.Unmarshal([]byte(data), &transaction)
			if err2 != nil {
				panic(err2)
			}
			totalTransactions = append(totalTransactions, transaction...)

		}
		totalVotes = len(totalTransactions)

		//get total presidential candidates
		candidates := []models.Candidate{}
		if err := database.Find(&candidates).Error; err != nil {
			log.Fatalln(err)
		}
		totalPresidentialCandidates = len(candidates)

		//get total desktop clients
		clients := []models.DesktopClient{}
		if err := database.Find(&clients).Error; err != nil {
			log.Fatalln(err)
		}
		totalDesktopClients = len(clients)

		//get total online desktop clients:todo

		//get transaction pool size
		transactions := []models.Transaction{}
		if err := database.Find(&transactions).Error; err != nil {
			log.Fatalln(err)
		}
		transactionPoolSize = len(transactions)

		//total processed transactions
		totalProcessedTransactions = totalVotes
	}
	data := models.QuickStats{
		TotalRegisteredVoters:      totalVoters,
		TotalPollingStations:       totalPollingStations,
		TotalVotes:                 totalVotes,
		PresidentialCandidates:     totalPresidentialCandidates,
		TotalDesktopClients:        totalDesktopClients,
		OnlineClients:              totalOnlineClients,
		TransactionPoolSize:        transactionPoolSize,
		TotalProcessedTransactions: totalProcessedTransactions,
	}
	context.IndentedJSON(http.StatusOK, data)

}
