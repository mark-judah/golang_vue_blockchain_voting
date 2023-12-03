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
	var newCounty []models.County
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
		for _, x := range newCounty {

			result := database.Omit(clause.Associations).Create(&x)
			if result.Error != nil {
				context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
			} else {
				context.IndentedJSON(http.StatusCreated, x)
				newAdminDashLog := models.AdminDashLog{Type: "County", Payload: x}
				PersistAdminDashLog(newAdminDashLog)

				mqttMessage := models.Message{Type: "new_county", Payload: x}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := Client[0].Publish("adminTransaction/1", 0, false, data)
				token.Wait()
			}
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
	var newConstituency []models.Constituency
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
		for _, x := range newConstituency {
			result := database.Create(&x)
			if result.Error != nil {
				context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
			} else {
				context.IndentedJSON(http.StatusCreated, x)
				newAdminDashLog := models.AdminDashLog{Type: "Constituency", Payload: x}
				PersistAdminDashLog(newAdminDashLog)
				mqttMessage := models.Message{Type: "new_constituency", Payload: x}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := Client[0].Publish("adminTransaction/1", 0, false, data)
				token.Wait()
			}
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
			ConstituencyID uint
			County         string
			Constituency   string
		}
		var constituencies []data

		for _, a := range county {
			for _, b := range a.Constituency {
				constituencies = append(constituencies, data{County: a.Name, Constituency: b.Name, ConstituencyID: b.ID})
			}
		}
		context.IndentedJSON(http.StatusOK, constituencies)
	}
}

func NewWard(context *gin.Context) {
	var newWard []models.Ward
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
		for _, x := range newWard {
			result := database.Create(&x)
			if result.Error != nil {
				context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
			} else {
				context.IndentedJSON(http.StatusCreated, x)
				newAdminDashLog := models.AdminDashLog{Type: "Ward", Payload: x}
				PersistAdminDashLog(newAdminDashLog)
				mqttMessage := models.Message{Type: "new_ward", Payload: x}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := Client[0].Publish("adminTransaction/1", 0, false, data)
				token.Wait()
			}
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
			WardID       uint
			County       string
			Constituency string
			Ward         string
		}
		var wards []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					wards = append(wards, data{County: a.Name, Constituency: b.Name, Ward: c.Name, WardID: c.ID})
				}
			}
		}
		context.IndentedJSON(http.StatusOK, wards)
	}
}

func NewPollingStation(context *gin.Context) {
	var newPollingStation []models.PollingStation
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
		for _, x := range newPollingStation {
			result := database.Create(&x)
			if result.Error != nil {
				context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
			} else {
				context.IndentedJSON(http.StatusCreated, x)
				newAdminDashLog := models.AdminDashLog{Type: "PollingStation", Payload: x}
				PersistAdminDashLog(newAdminDashLog)
				mqttMessage := models.Message{Type: "new_polling_station", Payload: x}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := Client[0].Publish("adminTransaction/1", 0, false, data)
				token.Wait()
			}
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
			PollingStationID uint
			County           string
			Constituency     string
			Ward             string
			PollingStation   string
		}
		var polling_stations []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					for _, d := range c.PollingStation {
						polling_stations = append(polling_stations, data{
							PollingStationID: d.ID,
							County:           a.Name,
							Constituency:     b.Name,
							Ward:             c.Name, PollingStation: d.Name,
						})
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
		fmt.Println(county)

		type data struct {
			ClientID       uint
			County         string
			Constituency   string
			Ward           string
			PollingStation string
			SerialNumber   string
		}
		var desktop_clients []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					for _, d := range c.PollingStation {
						for _, e := range d.DesktopClient {
							desktop_clients = append(desktop_clients, data{
								ClientID:       e.ID,
								County:         a.Name,
								Constituency:   b.Name,
								Ward:           c.Name,
								PollingStation: d.Name,
								SerialNumber:   e.SerialNumber})

						}
					}
				}
			}
		}
		context.IndentedJSON(http.StatusOK, desktop_clients)
	}
}

func NewCandidate(context *gin.Context) {
	var newCandidate []models.Candidate
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
		for _, x := range newCandidate {
			result := database.Create(&x)
			if result.Error != nil {
				context.IndentedJSON(http.StatusBadRequest, result.Error.Error())
			} else {
				context.IndentedJSON(http.StatusCreated, x)
				newAdminDashLog := models.AdminDashLog{Type: "Candidate", Payload: x}
				PersistAdminDashLog(newAdminDashLog)
				mqttMessage := models.Message{Type: "new_candidate", Payload: x}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := Client[0].Publish("adminTransaction/1", 0, false, data)
				token.Wait()
			}
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
			CandidateId    uint
			County         string
			Constituency   string
			Ward           string
			PollingStation string
			Candidate      string
			Position       string
			Party          string
			Slogan         string
			Statement      string
			Photo          string
		}
		var candidates []data

		for _, a := range county {
			for _, b := range a.Constituency {
				for _, c := range b.Ward {
					for _, d := range c.PollingStation {
						for _, e := range d.Candidate {
							candidates = append(candidates, data{
								CandidateId:    e.ID,
								County:         a.Name,
								Constituency:   b.Name,
								Ward:           c.Name,
								PollingStation: d.Name,
								Candidate:      e.Name,
								Position:       e.Position,
								Party:          e.Party,
								Slogan:         e.Slogan,
								Statement:      e.Statement,
								Photo:          e.Photo,
							})
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

func FetchRegions(context *gin.Context) {
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	counties := []models.County{}
	if err := database.Find(&counties).Error; err != nil {
		log.Fatalln(err)
	}
	constituencies := []models.Constituency{}
	if err := database.Find(&constituencies).Error; err != nil {
		log.Fatalln(err)
	}
	wards := []models.Ward{}
	if err := database.Find(&wards).Error; err != nil {
		log.Fatalln(err)
	}
	polling_stations := []models.PollingStation{}
	if err := database.Find(&polling_stations).Error; err != nil {
		log.Fatalln(err)
	}

	regions := models.Region{Counties: counties, Constituencies: constituencies, Wards: wards, PollingStations: polling_stations}
	context.IndentedJSON(http.StatusOK, regions)
}

func FetchTransactions(context *gin.Context) {
	database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	transactions := []models.Transaction{}
	if err := database.Find(&transactions).Error; err != nil {
		log.Fatalln(err)
	}

	type data struct {
		Txid           string
		NodeId         string
		Candidate      string
		County         string
		Constituency   string
		Ward           string
		PollingStation string
	}

	var transactionData []data

	for _, x := range transactions {
		candidate := models.Candidate{}
		county := models.County{}
		constituency := models.Constituency{}
		ward := models.Ward{}
		polling_station := models.PollingStation{}

		if err := database.Where("ID=?", x.CandidateId).First(&candidate).Error; err != nil {
			log.Fatalln(err)
		}
		if err := database.Where("ID=?", x.CountyID).First(&county).Error; err != nil {
			log.Fatalln(err)
		}
		if err := database.Where("ID=?", x.ConstituencyID).First(&constituency).Error; err != nil {
			log.Fatalln(err)
		}
		if err := database.Where("ID=?", x.WardID).First(&ward).Error; err != nil {
			log.Fatalln(err)
		}
		if err := database.Where("ID=?", x.PollingStationID).First(&polling_station).Error; err != nil {
			log.Fatalln(err)
		}

		transactionData = append(transactionData, data{
			Txid:           x.Txid,
			NodeId:         x.NodeId,
			Candidate:      candidate.Name,
			County:         county.Name,
			Constituency:   constituency.Name,
			Ward:           ward.Name,
			PollingStation: polling_station.Name,
		})
	}

	context.IndentedJSON(http.StatusOK, transactionData)

}
