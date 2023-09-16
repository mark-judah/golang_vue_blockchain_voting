package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
	"vote_backend/controller"
	"vote_backend/utils"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var leaderAlive = false
var leaderAliveCounter = 0
var leaderPayload []string
var myVotes []string

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var ctx = context.Background()

func main() {

	setRaftState("follower")
	setRaftTerm(0)
	setVoteAndTerm("0", "0", "0")

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")

	// interval to ping whether connectons to broker is still alive
	opts.SetKeepAlive(60 * time.Second)

	//set the handler to be called to receive messages when no subscription matches
	opts.SetDefaultPublishHandler(receiveMsgs)
	opts.SetPingTimeout(1 * time.Second)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	//create a new client, and pass the opts config to it
	//token.Wait, Wait() is a bool that show when a action is completed
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		leaderAlive = false

		key := readClientID() + "state"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("\n My state----->" + val)
		if !leaderAlive {
			leaderAliveCounter = leaderAliveCounter + 1
		}
		if leaderAliveCounter >= 10 && val == "follower" {
			fmt.Println("Leader Dead")
			requestVotes(client)

		}

		if val == "leader" {
			key := readClientID() + "term"
			val, err := redisClient.Get(ctx, key).Result()
			if err != nil {
				panic(err)
			}

			myPayload := append(leaderPayload, val, "Leader Alive")

			jsonData, err2 := json.Marshal(myPayload)
			if err2 != nil {
				panic(err2)
			}

			time.Sleep(time.Duration(time.Second))
			// publish a message every one second
			token := client.Publish("leaderNodePulse/"+key, 0, false, jsonData)
			token.Wait()

			//check transaction pool
			println("Transaction pool: " + fmt.Sprintf("%+v", controller.TransactionPool.Transactions))
			transactionData, err3 := json.Marshal(controller.TransactionPool)
			if err3 != nil {
				panic(err3)
			}
			token2 := client.Publish("transactionPool/1", 0, false, transactionData)
			token2.Wait()
		}

		if val == "candidate" {
			//set to zero to avoid requestVotes being called infinitely
			leaderAliveCounter = 0
			go killGinServer()

		}

		if val == "follower" {
			fmt.Println("\n Leader Alive--------------------->" + strconv.FormatBool(leaderAlive))
			fmt.Println("\n Leader Alive Counter--------------------->" + strconv.Itoa(leaderAliveCounter))
			fmt.Println(getClientState())
			go killGinServer()
		}

		time.Sleep(time.Duration(time.Second))

		fmt.Println("\n My id----->" + readClientID())

	}

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

	if token := client.Subscribe("transactionPool/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection to broker lost: %v", err)

}
var receiveMsgs mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	var dataArray []string
	var transactionPool utils.Queue
	err := json.Unmarshal(message.Payload(), &dataArray)
	if err != nil {
		err := json.Unmarshal(message.Payload(), &transactionPool)
		println("Transaction pool message: " + fmt.Sprintf("%+v", transactionPool.Transactions))
		if err != nil {
			log.Println("Error unmarshalling as string or as Queue:", err)
		}

	}
	fmt.Printf("TOPIC------> %s\n", message.Topic())
	fmt.Printf("MESSAGE------> %s\n", string(message.Payload()))

	if len(dataArray) > 0 {
		if dataArray[1] == "Leader Alive" {
			println("-------------------->Leader Alive")
			if getClientState() != "leader" {
				setRaftState("follower")
			}
			time.Sleep(time.Duration(time.Second))
			leaderAlive = true
			leaderAliveCounter = 0
		}
	}

	if string(message.Topic()) == "election/1" {
		fmt.Println("Node:" + readClientID() + " casting vote")
		intVal, err2 := strconv.Atoi(dataArray[1])
		if err2 != nil {
			panic(err2)
		}

		if getClientTerm() > intVal {
			fmt.Println("Node has a greater term than candidate")
			setRaftState("candidate")
			candidateNodeId := dataArray[0]
			var voterPayload []string
			voterPayload = append(voterPayload, readClientID(), candidateNodeId, strconv.Itoa(intVal), "higher term")
			jsonData, err2 := json.Marshal(voterPayload)
			if err2 != nil {
				panic(err2)
			}
			token := client.Publish("nodeElectionResponse/1", 0, false, jsonData)
			token.Wait()
		}

		//check if the node has already voted in this election term

		if getClientVote()[1] != dataArray[1] {
			fmt.Println("Node has not voted in this term" + " stored candidate: " + getClientVote()[0] + " election term " + dataArray[1])

			if getClientState() == "candidate" {
				println("Node voting for itself")
				var voterPayload []string
				voterPayload = append(voterPayload, readClientID(), readClientID(), strconv.Itoa(intVal), "yes")
				jsonData, err2 := json.Marshal(voterPayload)
				if err2 != nil {
					panic(err2)
				}
				token := client.Publish("nodeElectionResponse/1", 0, false, jsonData)
				token.Wait()

				setVoteAndTerm(readClientID(), voterPayload[2], "yes")
			} else {
				println("Node voting for another node")
				candidateNodeId := dataArray[0]
				var voterPayload []string
				voterPayload = append(voterPayload, readClientID(), candidateNodeId, strconv.Itoa(intVal), "yes")
				jsonData, err2 := json.Marshal(voterPayload)
				if err2 != nil {
					panic(err2)
				}
				token := client.Publish("nodeElectionResponse/1", 0, false, jsonData)
				token.Wait()

				setVoteAndTerm(candidateNodeId, voterPayload[2], "yes")
			}
		} else {
			fmt.Println("Node has already voted in this term" + " stored candidate: " + getClientVote()[0] + " election term" + dataArray[1])
		}

	}

	if string(message.Topic()) == "nodeElectionResponse/1" {
		fmt.Println("Node:" + readClientID() + "vote response")

		if readClientID() == dataArray[1] {
			if dataArray[3] == "higher term" {
				setRaftState("follower")
			}

			if len(myVotes) < getTotalConnectedNodes() {
				myVotes = append(myVotes, dataArray[1])
			}

			if len(myVotes) == getTotalConnectedNodes() {
				setRaftState("leader")
				go startGinServer()
			}

		}

	}

	if string(message.Topic()) == "transactionPool/1" {
		fmt.Println("\n --------------------->" + "Transactions")
		//verify the transactions

	}
}

var requestVotes = func(client mqtt.Client) {
	fmt.Println("\n --------------------->" + "Voting Started")
	randomNumber := rand.Intn(10)
	fmt.Println("\n --------------------->" + "Waited for" + strconv.Itoa(randomNumber))
	time.Sleep(time.Duration(time.Duration(randomNumber).Seconds()))
	setRaftState("candidate")

	key := readClientID() + "term"
	term, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}

	newTerm, err2 := strconv.Atoi(term)
	if err != nil {
		panic(err2)
	}
	newTerm = newTerm + 1
	setRaftTerm(newTerm)

	var candidatePayload []string
	candidatePayload = append(candidatePayload, readClientID(), strconv.Itoa(newTerm))
	jsonData, err2 := json.Marshal(candidatePayload)
	if err2 != nil {
		panic(err2)
	}
	token := client.Publish("election/1", 0, false, jsonData)
	token.Wait()
}

func getClientTerm() int {
	key := readClientID() + "term"
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	intVal, err2 := strconv.Atoi(val)
	if err2 != nil {
		panic(err2)
	}
	return intVal
}

func getClientState() string {
	key := readClientID() + "state"
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func getClientVote() []string {
	key := readClientID() + "votePayload"
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	var dataArray []string
	err = json.Unmarshal([]byte(val), &dataArray)
	if err != nil {
		log.Println("Error unmarshalling:", err)
	}

	return dataArray
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

func readClientID() string {
	if _, err := os.Stat("clientId"); err == nil {
		fmt.Printf("File exists\n")

		clientId, err := os.ReadFile("clientId")
		if err != nil {
			panic(err)
		}
		return string(clientId)
	} else {
		fmt.Printf("File does not exist\n")

		clientId := []byte(uuid.New().String())
		err := os.WriteFile("clientId", clientId, 0644)
		if err != nil {
			panic(err)
		}
		return string(clientId)
	}

}

// func readClientID() string {
// 	key := "clientID"
// 	val, err := redisClient.Get(ctx, key).Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return val
// }

func setRaftTerm(term int) {
	err := redisClient.Set(ctx, readClientID()+"term", term, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := readClientID() + "term"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Term", val)
	}
}

func setRaftState(state string) {
	err := redisClient.Set(ctx, readClientID()+"state", state, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := readClientID() + "state"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("State", val)
	}
}

func setVoteAndTerm(candidateNodeId string, term string, vote string) {
	fmt.Println("Storing Vote: " + "candidateNodeId: " + candidateNodeId + " term: " + term + " vote: " + vote)
	var votePayload []string
	votePayload = append(votePayload, candidateNodeId, term, vote)
	jsonData, err2 := json.Marshal(votePayload)
	if err2 != nil {
		panic(err2)
	}

	err := redisClient.Set(ctx, readClientID()+"votePayload", jsonData, 0).Err()
	if err != nil {
		panic(err)
	} else {
		val, err := redisClient.Get(ctx, readClientID()+"votePayload").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Stored vote", val)
	}
}

func startGinServer() {
	//only the leader can create a router and receive requests
	router := gin.Default()
	router.POST("/new-vote", controller.NewVote)
	router.Run("localhost:3500")
}

func killGinServer() {
	command := fmt.Sprintf("lsof -i tcp:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", "3500")
	exec_cmd(exec.Command("/bin/bash", "-c", command))
}

func exec_cmd(cmd *exec.Cmd) {
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("Error during killing (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Port successfully killed (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}
