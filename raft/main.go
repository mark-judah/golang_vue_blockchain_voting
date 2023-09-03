package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var leaderAlive = false
var leaderAliveCounter = 0

var leaderPayload []string

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var ctx = context.Background()

func main() {

	err := redisClient.Set(ctx, getClientID()+"nodeID", getClientID(), 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := getClientID() + "nodeID"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("nodeID", val)
	}

	setRaftState("leader")
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

		key := getClientID() + "state"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}

		if !leaderAlive {
			leaderAliveCounter = leaderAliveCounter + 1
		}
		if leaderAliveCounter >= 10 && val == "follower" {
			fmt.Println("Leader Dead")
			requestVotes(client)

		}

		if val == "leader" {
			key := getClientID() + "term"
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

		}

		if val == "candidate" {
			//set to zero to avoid requestVotes being called infinitely
			leaderAliveCounter = 0

		}
		time.Sleep(time.Duration(time.Second))
		fmt.Print(leaderAlive)
		fmt.Print(leaderAliveCounter)
		fmt.Print(getClientState())
		getTotalNodes()
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
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect to broker lost: %v", err)
}
var receiveMsgs mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	var dataArray []string
	err := json.Unmarshal(message.Payload(), &dataArray)
	if err != nil {
		log.Println("Error unmarshalling:", err)
		return
	}
	fmt.Println("NEW UNMARSHALLED MESSAGE: " + dataArray[0] + " " + dataArray[1] + " ")

	fmt.Printf("TOPIC: %s\n", message.Topic())
	fmt.Printf("MESSAGE: %s\n", string(message.Payload()))

	if dataArray[1] == "Leader Alive" {
		print("##################################")
		if getClientState() != "leader" {
			setRaftState("follower")
		}
		time.Sleep(time.Duration(time.Second))
		leaderAlive = true
		leaderAliveCounter = 0
	}

	if string(message.Topic()) == "election/1" {
		fmt.Println("Node:" + getClientID() + " casting vote")
		intVal, err2 := strconv.Atoi(dataArray[1])
		if err2 != nil {
			panic(err2)
		}

		if getClientTerm() > intVal {
			fmt.Println("Node has a greater term than candidate")
			setRaftState("candidate")
			candidateNodeId := dataArray[0]
			var voterPayload []string
			voterPayload = append(voterPayload, getClientID(), candidateNodeId, strconv.Itoa(intVal), "higher term")
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
				voterPayload = append(voterPayload, getClientID(), getClientID(), strconv.Itoa(intVal), "yes")
				jsonData, err2 := json.Marshal(voterPayload)
				if err2 != nil {
					panic(err2)
				}
				token := client.Publish("nodeElectionResponse/1", 0, false, jsonData)
				token.Wait()

				setVoteAndTerm(getClientID(), voterPayload[2], "yes")
			} else {
				println("Node voting for another node")
				candidateNodeId := dataArray[0]
				var voterPayload []string
				voterPayload = append(voterPayload, getClientID(), candidateNodeId, strconv.Itoa(intVal), "yes")
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
		fmt.Println("Node:" + getClientID() + "vote response")

		if getClientID() == dataArray[1] {
			if dataArray[3] == "higher term" {
				setRaftState("follower")
			}
		}

		//store each vote as it arrives
		votesMap := map[string]string{}
		votesMap[dataArray[0]] = dataArray[2]

		// //count total number of votes
		// if(len(votesMap)==3){
		// 	for candidate,vote:=range votesMap{

		// 	}
		// }
		//announce leader
		//revert all other nodes to followers
	}

}

var requestVotes = func(client mqtt.Client) {
	fmt.Println("########################" + "Voting Started")
	randomNumber := rand.Intn(10)
	fmt.Println("########################" + "Waited for" + strconv.Itoa(randomNumber))
	time.Sleep(time.Duration(time.Duration(randomNumber).Seconds()))
	setRaftState("candidate")

	key := getClientID() + "term"
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
	candidatePayload = append(candidatePayload, getClientID(), strconv.Itoa(newTerm))
	jsonData, err2 := json.Marshal(candidatePayload)
	if err2 != nil {
		panic(err2)
	}
	token := client.Publish("election/1", 0, false, jsonData)
	token.Wait()
}

func getClientTerm() int {
	key := getClientID() + "term"
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
	key := getClientID() + "state"
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func getClientVote() []string {
	key := getClientID() + "votePayload"
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

func getTotalNodes() {

}

func getClientID() string {
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

func setRaftTerm(term int) {
	err := redisClient.Set(ctx, getClientID()+"term", term, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := getClientID() + "term"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Term", val)
	}
}

func setRaftState(state string) {
	err := redisClient.Set(ctx, getClientID()+"state", state, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := getClientID() + "state"
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

	err := redisClient.Set(ctx, getClientID()+"votePayload", jsonData, 0).Err()
	if err != nil {
		panic(err)
	} else {
		val, err := redisClient.Get(ctx, getClientID()+"votePayload").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Stored vote", val)
	}
}
