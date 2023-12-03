package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
	"vote_backend/models"
	utils "vote_backend/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

var LeaderAlive = false
var LeaderAliveCounter = 0
var LeaderPayload []string
var MyVotes []string
var Client []mqtt.Client

func InitMqttClient() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")

	// interval to ping whether connectons to broker is still alive
	opts.SetKeepAlive(60 * time.Second)

	// Set the handler to be called to receive messages when no subscription matches
	opts.SetDefaultPublishHandler(receiveMsgs)
	opts.SetPingTimeout(1 * time.Second)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// create a new client, and pass the opts config to it
	// token.Wait, Wait() is a bool that show when a action is completed
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	Client = append(Client, client)

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
	//subcribe to a topic
	//token.Wait, Wait() is a bool that show when a action is completed

	//place subscriptions here
	//on connection loss, the client resubscribes when the connection is restored
	if token := client.Subscribe("leaderNodePulse/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("election/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("nodeElectionResponse/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("followerAppend/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("raftLogAppendConfirmation/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("leaderLogRequest/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("leaderLogResponse/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("tallyVotes/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("tallyResults/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	if token := client.Subscribe("adminTransaction/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	if token := client.Subscribe("leaderAdminDashLogRequest/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	if token := client.Subscribe("leaderAdminDashLogResponse/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	if token := client.Subscribe("nodeStatsRequest/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	if token := client.Subscribe("nodeStatsResponse/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection to broker lost: %v", err)

}
var receiveMsgs mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	var leaderTermPulse []string
	var nodeElection []string
	var nodeElectionResponse []string

	var newTransaction models.Transaction
	var leaderUser models.Users
	var leaderCounty models.County
	var leaderConstituency models.Constituency
	var leaderWard models.Ward
	var leaderPollingStation models.PollingStation
	var leaderDesktopClient models.DesktopClient
	var leaderCandidate models.Candidate
	var leaderVoter models.Voter

	var leaderTransactions []models.Transaction
	var tallyResults map[string]int

	var leaderAdminDashLogResponse []models.AdminDashLog

	var newNodeStats models.NodeStats

	var msg models.Message
	json.Unmarshal(message.Payload(), &msg)

	switch msg.Type {
	case "leader_term_pulse":
		fmt.Println("type leader_term_pulse")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderTermPulse)

	case "node_election":
		fmt.Println("type node_election")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &nodeElection)

	case "node_election_response":
		fmt.Println("type node_election_response")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &nodeElectionResponse)

	case "new_user":
		fmt.Println("type new_user")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderUser)

	case "new_county":
		fmt.Println("type new_county")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderCounty)

	case "new_constituency":
		fmt.Println("type new_constituency")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderConstituency)

	case "new_ward":
		fmt.Println("type new_ward")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderWard)

	case "new_polling_station":
		fmt.Println("type new_polling_station")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderPollingStation)

	case "new_desktop_client":
		fmt.Println("type new_desktop_client")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderDesktopClient)

	case "new_candidate":
		fmt.Println("type new_candidate")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderCandidate)

	case "new_voter":
		fmt.Println("type new_voter")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderVoter)

	case "new_transaction":
		fmt.Println("type new_transaction")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &newTransaction)

	case "leader_log_response":
		fmt.Println("type leader_log_response")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderTransactions)

	case "tally_results":
		fmt.Println("type tally_results")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &tallyResults)

	case "leader_admin_dash_log_response":
		fmt.Println("type leader_admin_dash_log_response")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &leaderAdminDashLogResponse)

	case "node_stats_response":
		fmt.Println("type node_stats_response")
		jsonBytes, _ := json.Marshal(msg.Payload)
		json.Unmarshal(jsonBytes, &newNodeStats)
	}

	fmt.Printf("TOPIC------> %s\n", message.Topic())
	fmt.Printf("MESSAGE------> %s\n", string(message.Payload()))

	if len(leaderTermPulse) > 0 {
		if leaderTermPulse[1] == "Leader Alive" {
			println("-------------------->Leader Alive")

			if utils.GetClientState() != "leader" {
				if utils.GetClientState() == "syncing" {
					fmt.Println("Skipping syncing node")
				} else {
					utils.SetRaftState("follower")
				}
			}
			time.Sleep(time.Duration(time.Second))
			LeaderAlive = true
			LeaderAliveCounter = 0
		}
	}

	if string(message.Topic()) == "election/1" {
		fmt.Println("Node:" + utils.ReadClientID() + " casting vote")
		candidateTerm, err2 := strconv.Atoi(nodeElection[1])
		if err2 != nil {
			panic(err2)
		}

		if utils.GetClientTerm() > candidateTerm {
			fmt.Println("Node has a greater term than candidate")
			utils.SetRaftState("candidate")
			candidateNodeId := nodeElection[0]
			var voterPayload []string
			voterPayload = append(voterPayload, utils.ReadClientID(), candidateNodeId, strconv.Itoa(candidateTerm), "higher term")

			mqttMessage := models.Message{Type: "node_election_response", Payload: voterPayload}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := client.Publish("nodeElectionResponse/1", 0, false, data)
			token.Wait()
		}

		//check if the node has already voted in this election term

		if utils.GetClientVote()[1] != nodeElection[1] {
			fmt.Println("Node has not voted in this term" + " stored candidate: " + utils.GetClientVote()[0] + " election term " + nodeElection[1])

			if utils.GetClientState() == "candidate" {
				println("Node voting for itself")
				var voterPayload []string
				voterPayload = append(voterPayload, utils.ReadClientID(), utils.ReadClientID(), strconv.Itoa(candidateTerm), "yes")
				mqttMessage := models.Message{Type: "node_election_response", Payload: voterPayload}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := client.Publish("nodeElectionResponse/1", 0, false, data)
				token.Wait()

				utils.SetVoteAndTerm(utils.ReadClientID(), voterPayload[2], "yes")
			} else {
				println("Node voting for another node")
				candidateNodeId := nodeElection[0]
				var voterPayload []string
				voterPayload = append(voterPayload, utils.ReadClientID(), candidateNodeId, strconv.Itoa(candidateTerm), "yes")
				mqttMessage := models.Message{Type: "node_election_response", Payload: voterPayload}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := client.Publish("nodeElectionResponse/1", 0, false, data)
				token.Wait()

				utils.SetVoteAndTerm(candidateNodeId, voterPayload[2], "yes")
			}
		} else {
			fmt.Println("Node has alReady voted in this term" + " stored candidate: " + utils.GetClientVote()[0] + " election term" + nodeElection[1])
		}

	}

	if string(message.Topic()) == "nodeElectionResponse/1" {
		fmt.Println("Node:" + utils.ReadClientID() + "vote response")

		if utils.ReadClientID() == nodeElectionResponse[1] {
			if nodeElectionResponse[3] == "higher term" {
				utils.SetRaftState("follower")
			}

			if len(MyVotes) < getTotalConnectedNodes() {
				MyVotes = append(MyVotes, nodeElectionResponse[1])
			}

			//todo: change to half total nodes due to latency
			if len(MyVotes) == getTotalConnectedNodes() {
				utils.SetRaftState("leader")
				go StartApiServer()
			}

		}

	}

	if string(message.Topic()) == "followerAppend/1" {
		fmt.Println("\n --------------------->" + "New Transaction to Append")

		//update log with new transaction
		if utils.GetClientState() == "follower" {
			fmt.Println("Writing to db" + fmt.Sprintf("%+v", newTransaction))
			database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
			if err != nil {
				panic(err)
			}

			//commit to log.json file
			PersistLog(newTransaction)

			// if ValidateVote(newTransaction) {
			database.Create(&newTransaction)
			//}

			database.Create(&newTransaction)

		}
	}

	if string(message.Topic()) == "leaderLogRequest/1" {
		var transactions []models.Transaction

		if utils.GetClientState() == "leader" {
			if _, err := os.Stat("log.json"); errors.Is(err, os.ErrNotExist) {
				//log.json doesn't exist,retur empty transactions
				mqttMessage := models.Message{Type: "leader_log_response", Payload: transactions}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := client.Publish("leaderLogResponse/1", 0, false, data)
				token.Wait()
				return
			}
			logFile, err := os.ReadFile("log.json")
			if err != nil {
				fmt.Println(err)
			}

			err2 := json.Unmarshal(logFile, &transactions)
			if err2 != nil {
				panic(err2)
			}
			mqttMessage := models.Message{Type: "leader_log_response", Payload: transactions}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := client.Publish("leaderLogResponse/1", 0, false, data)
			token.Wait()

		}
	}

	if string(message.Topic()) == "leaderLogResponse/1" {
		if utils.GetClientState() == "syncing" {
			//check if leaderLog is empty
			if len(leaderTransactions) == 0 {
				fmt.Println("The leader has no transactions to sync yet")
				SyncAdminDashLog()
				return
			}
			//check if file exists
			_, err := os.Stat("log.json")
			if err != nil {
				//log file doesn't exist
				for _, x := range leaderTransactions {
					fmt.Println("Writing to db" + fmt.Sprintf("%+v", x))
					database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
					if err != nil {
						panic(err)
					}
					//commit to log.json file
					PersistLog(x)

					//if ValidateVote(x) {
					database.Create(&x)
					//}

				}
				SyncAdminDashLog()

			} else {
				//if file exists, check if follower log entries match the leadernodes entries, if not append from the last matching value
				syncingNodesLog, err := os.ReadFile("log.json")
				if err != nil {
					panic(err)
				}
				var followerTransactions []models.Transaction

				err2 := json.Unmarshal(syncingNodesLog, &followerTransactions)
				if err2 != nil {
					panic(err2)
				}
				if len(followerTransactions) < len(leaderTransactions) {
					fmt.Println("Mismatch, leaderlog length:" + strconv.Itoa(len(leaderTransactions)))
					fmt.Println("Mismatch, follower log length:" + strconv.Itoa(len(followerTransactions)))

					for i := len(followerTransactions); i < len(leaderTransactions); i++ {
						PersistLog(leaderTransactions[i])
					}
				}
				//as soon as the leadernode file size matches the syncing node file size, set the raft state to follower,
				// if it doesn't match call the nodeSync function recursively
				//set state to follower once it is confirmed that the logs match

				if len(followerTransactions) == len(leaderTransactions) {
					fmt.Println("Logs match")
					SyncAdminDashLog()

				} else {
					fmt.Println("leaderlog length:" + strconv.Itoa(len(leaderTransactions)))
					fmt.Println("follower log length:" + strconv.Itoa(len(followerTransactions)))

					NodeSync()
				}
			}
		}

	}

	if string(message.Topic()) == "tallyVotes/1" {
		TallyFeedback()
	}

	if string(message.Topic()) == "tallyResults/1" {
		fmt.Println(tallyResults)
		fmt.Println("pass tally")
		FinalTally(tallyResults)
	}
	if string(message.Topic()) == "adminTransaction/1" {
		if utils.GetClientState() == "follower" {

			database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
			if err != nil {
				panic(err)
			}
			switch {
			case !reflect.ValueOf(leaderUser).IsZero():
				result := database.Create(&leaderUser)
				if result.Error != nil {
					panic(result.Error)
				}
			case !reflect.ValueOf(leaderCounty).IsZero():
				result := database.Create(&leaderCounty)
				if result.Error != nil {
					panic(result.Error)
				}
			case reflect.ValueOf(leaderConstituency).IsZero():
				result := database.Create(&leaderConstituency)
				if result.Error != nil {
					panic(result.Error)
				}
			case reflect.ValueOf(leaderWard).IsZero():
				result := database.Create(&leaderWard)
				if result.Error != nil {
					panic(result.Error)
				}
			case reflect.ValueOf(leaderPollingStation).IsZero():
				result := database.Create(&leaderPollingStation)
				if result.Error != nil {
					panic(result.Error)
				}
			case reflect.ValueOf(leaderDesktopClient).IsZero():
				result := database.Create(&leaderDesktopClient)
				if result.Error != nil {
					panic(result.Error)
				}
			case reflect.ValueOf(leaderCandidate).IsZero():
				result := database.Create(&leaderCandidate)
				if result.Error != nil {
					panic(result.Error)
				}
			case reflect.ValueOf(leaderVoter).IsZero():
				result := database.Create(&leaderVoter)
				if result.Error != nil {
					panic(result.Error)
				}
			}
		}
	}
	if string(message.Topic()) == "leaderAdminDashLogRequest/1" {
		var admin_dash_log []models.AdminDashLog

		if utils.GetClientState() == "leader" {
			if _, err := os.Stat("admin_dash_log.json"); errors.Is(err, os.ErrNotExist) {
				//log.json doesn't exist,retur empty transactions
				mqttMessage := models.Message{Type: "leader_admin_dash_log_response", Payload: admin_dash_log}
				data, err := json.Marshal(mqttMessage)
				if err != nil {
					panic(err)
				}
				token := client.Publish("leaderLogResponse/1", 0, false, data)
				token.Wait()
				return
			}
			logFile, err := os.ReadFile("admin_dash_log.json")
			if err != nil {
				fmt.Println(err)
			}

			err2 := json.Unmarshal(logFile, &admin_dash_log)
			if err2 != nil {
				panic(err2)
			}
			mqttMessage := models.Message{Type: "leader_admin_dash_log_response", Payload: admin_dash_log}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}
			token := client.Publish("leaderAdminDashLogResponse/1", 0, false, data)
			token.Wait()

		}
	}

	if string(message.Topic()) == "leaderAdminDashLogResponse/1" {
		if utils.GetClientState() == "syncing" {

			//delete old log file
			err := os.Remove("admin_dash_log.json")
			if err != nil {
				fmt.Println("dmin dash log ile doesn't exist")
			}
			database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
			if err != nil {
				panic(err)
			}
			for _, x := range leaderAdminDashLogResponse {
				fmt.Println(x)
				switch x.Type {
				case "User":
					var user_model models.Users
					user, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(user, &user_model)
					result := database.Create(&user_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				case "County":
					var couny_model models.County
					county, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(county, &couny_model)
					result := database.Create(&couny_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				case "Constituency":
					var constituency_model models.Constituency
					constituency, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(constituency, &constituency_model)
					result := database.Create(&constituency_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				case "Ward":
					var ward_model models.Ward
					ward, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(ward, &ward_model)
					result := database.Create(&ward_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				case "PollingStation":
					var polling_station_model models.PollingStation
					polling_station, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(polling_station, &polling_station_model)
					result := database.Create(&polling_station_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				case "DesktopClient":
					var desktop_client_model models.DesktopClient
					desktop_client, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(desktop_client, &desktop_client_model)
					result := database.Create(&desktop_client_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				case "Candidate":
					var candidate_model models.Candidate
					candidate, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(candidate, &candidate_model)
					result := database.Create(&candidate_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				case "Voter":
					var voter_model models.Voter
					voter, err := json.Marshal(x.Payload)
					if err != nil {
						panic(err)
					}
					json.Unmarshal(voter, &voter_model)
					result := database.Create(&voter_model)
					if result.Error != nil {
						fmt.Println(result.Error)
					}
					PersistAdminDashLog(x)

				}
			}

			//set raft state to follower
			utils.SetRaftState("follower")
		}
	}

	if string(message.Topic()) == "nodeStatsRequest/1" {
		var nodesTransactions []models.Transaction
		var logFileLength int
		logFile, err := os.ReadFile("log.json")
		if err != nil {
			fmt.Println("No log file yet")
			logFileLength = 0
		} else {
			err2 := json.Unmarshal(logFile, &nodesTransactions)
			if err2 != nil {
				panic(err2)
			}
			logFileLength = len(nodesTransactions)
		}

		node_stats := models.NodeStats{NodeId: utils.ReadClientID(), Status: utils.GetClientState(), Term: utils.GetClientTerm(), LogLength: logFileLength, DashboardLink: utils.GetClientPort()}
		mqttMessage := models.Message{Type: "node_stats_response", Payload: node_stats}
		data, err := json.Marshal(mqttMessage)
		if err != nil {
			panic(err)
		}
		token := Client[0].Publish("nodeStatsResponse/1", 0, false, data)
		token.Wait()
	}

	if string(message.Topic()) == "nodeStatsResponse/1" {
		NodeStats = append(NodeStats, newNodeStats)
	}

}

func getTotalConnectedNodes() int {
	username := "920ed6b2165341ff"
	password := "YXVq6EHQTA5dNyLcvvFiuQO6KfJT33wegV9B8MR4fz3C"

	url := "http://localhost:18083/api/v5/nodes/emqx%40127.0.0.1/stats"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := io.ReadAll(resp.Body)
	var data map[string]int
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		log.Println("\n Error unmarshalling:", err)
	}
	live_connections := data["live_connections.count"]
	fmt.Println("Total live connections: " + strconv.Itoa(live_connections))

	return live_connections
}
