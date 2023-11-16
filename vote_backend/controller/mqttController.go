package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection to broker lost: %v", err)

}
var receiveMsgs mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	var dataArray []string
	var newTransaction models.Transaction
	var responseString string
	var leaderTransactions []models.Transaction

	err1 := json.Unmarshal(message.Payload(), &dataArray)
	if err1 != nil {
		err2 := json.Unmarshal(message.Payload(), &newTransaction)
		println("Raft Log content: " + fmt.Sprintf("%+v", newTransaction))
		if err2 != nil {
			err3 := json.Unmarshal(message.Payload(), &responseString)
			println("Append Response: " + fmt.Sprintf("%+v", &responseString))
			if err3 != nil {
				err4 := json.Unmarshal(message.Payload(), &leaderTransactions)
				if err4 != nil {
					fmt.Println("Error unmarshalling as string,  string array Map or as transactions array:", err4)
				}
			}
		}

	}
	fmt.Printf("TOPIC------> %s\n", message.Topic())
	fmt.Printf("MESSAGE------> %s\n", string(message.Payload()))

	if len(dataArray) > 0 {
		if dataArray[1] == "Leader Alive" {
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
		intVal, err2 := strconv.Atoi(dataArray[1])
		if err2 != nil {
			panic(err2)
		}

		if utils.GetClientTerm() > intVal {
			fmt.Println("Node has a greater term than candidate")
			utils.SetRaftState("candidate")
			candidateNodeId := dataArray[0]
			var voterPayload []string
			voterPayload = append(voterPayload, utils.ReadClientID(), candidateNodeId, strconv.Itoa(intVal), "higher term")
			jsonData, err2 := json.Marshal(voterPayload)
			if err2 != nil {
				panic(err2)
			}
			token := client.Publish("nodeElectionResponse/1", 0, false, jsonData)
			token.Wait()
		}

		//check if the node has already voted in this election term

		if utils.GetClientVote()[1] != dataArray[1] {
			fmt.Println("Node has not voted in this term" + " stored candidate: " + utils.GetClientVote()[0] + " election term " + dataArray[1])

			if utils.GetClientState() == "candidate" {
				println("Node voting for itself")
				var voterPayload []string
				voterPayload = append(voterPayload, utils.ReadClientID(), utils.ReadClientID(), strconv.Itoa(intVal), "yes")
				jsonData, err2 := json.Marshal(voterPayload)
				if err2 != nil {
					panic(err2)
				}
				token := client.Publish("nodeElectionResponse/1", 0, false, jsonData)
				token.Wait()

				utils.SetVoteAndTerm(utils.ReadClientID(), voterPayload[2], "yes")
			} else {
				println("Node voting for another node")
				candidateNodeId := dataArray[0]
				var voterPayload []string
				voterPayload = append(voterPayload, utils.ReadClientID(), candidateNodeId, strconv.Itoa(intVal), "yes")
				jsonData, err2 := json.Marshal(voterPayload)
				if err2 != nil {
					panic(err2)
				}
				token := client.Publish("nodeElectionResponse/1", 0, false, jsonData)
				token.Wait()

				utils.SetVoteAndTerm(candidateNodeId, voterPayload[2], "yes")
			}
		} else {
			fmt.Println("Node has alReady voted in this term" + " stored candidate: " + utils.GetClientVote()[0] + " election term" + dataArray[1])
		}

	}

	if string(message.Topic()) == "nodeElectionResponse/1" {
		fmt.Println("Node:" + utils.ReadClientID() + "vote response")

		if utils.ReadClientID() == dataArray[1] {
			if dataArray[3] == "higher term" {
				utils.SetRaftState("follower")
			}

			if len(MyVotes) < getTotalConnectedNodes() {
				MyVotes = append(MyVotes, dataArray[1])
			}

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

			//verify the transactions...call  a vote verification function, return bool
			//ensure that the txid is valid, and was generated by an official client app
			//ensure that the node id is valid
			//ensure that the candidate id is valid
			//ensure that the vote hash is valid
			//ensure that the voter exists and that the voters details hash matches the stored hash
			//ensure that the voter hasnt already voted
			//insert verified transaction into db
			database.Create(&newTransaction)

		}
	}

	if string(message.Topic()) == "leaderLogRequest/1" {
		if utils.GetClientState() == "leader" {
			logFile, err := os.ReadFile("log.json")
			if err != nil {
				panic(err)
			}
			var transactions []models.Transaction

			err2 := json.Unmarshal(logFile, &transactions)
			if err2 != nil {
				panic(err2)
			}
			jsonData, err3 := json.Marshal(transactions)
			if err3 != nil {
				panic(err3)
			}
			token := client.Publish("leaderLogResponse/1", 0, false, jsonData)
			token.Wait()

		}
	}

	if string(message.Topic()) == "leaderLogResponse/1" {
		if utils.GetClientState() == "syncing" {
			for _, x := range leaderTransactions {
				fmt.Println("Writing to db" + fmt.Sprintf("%+v", x))
				database, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
				if err != nil {
					panic(err)
				}
				//commit to log.json file
				PersistLog(x)

				//verify the transactions...call  a vote verification function, return bool
				//ensure that the txid is valid, and was generated by an official client app
				//ensure that the node id is valid
				//ensure that the candidate id is valid
				//ensure that the vote hash is valid
				//ensure that the voter exists and that the voters details hash matches the stored hash
				//ensure that the voter hasnt already voted
				//insert verified transaction into db
				database.Create(&x)
			}
		}

		//if file exists, check if follower log entries match the leader nodes entries, if not append from the last matching value
		//as soon as the leader node file size matches the syncing node file size, set the raft state to leader, if it doen;t match
		//call the nodeSync function recursively
		//set state to follower once it is confirmed that the logs match
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
