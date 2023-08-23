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

var clientID = uuid.New()
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

	err := redisClient.Set(ctx, clientID.String()+"nodeID", clientID.String(), 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := clientID.String() + "nodeID"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("nodeID", val)
	}

	setRaftState("leader")
	setRaftTerm(0)

	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://broker.emqx.io:1883")

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

	//subcribe to a topic
	//token.Wait, Wait() is a bool that show when a action is completed

	if token := client.Subscribe("leaderNodePulse/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	if token := client.Subscribe("election/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		leaderAlive = false

		key := clientID.String() + "state"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}

		if !leaderAlive {
			leaderAliveCounter = leaderAliveCounter + 1
		}
		if leaderAliveCounter >= 10 && val == "follower" {
			fmt.Printf("Leader Dead")
			requestVotes(client)

		}

		if val == "leader" {
			key := clientID.String() + "term"
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
	}

}

func getClientState() string {
	key := clientID.String() + "state"
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

func setRaftTerm(term int) {
	err := redisClient.Set(ctx, clientID.String()+"term", term, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := clientID.String() + "term"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("Term", val)
	}
}

func setRaftState(state string) {
	err := redisClient.Set(ctx, clientID.String()+"state", state, 0).Err()
	if err != nil {
		panic(err)
	} else {
		key := clientID.String() + "state"
		val, err := redisClient.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("State", val)
	}
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}
var receiveMsgs mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	var dataArray []string
	err := json.Unmarshal(message.Payload(), &dataArray)
	if err != nil {
		log.Println("Error unmarshalling:", err)
		return
	}

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

	if string(message.Topic()) == "election" {
		fmt.Printf("Election")

	}
}

var requestVotes = func(client mqtt.Client) {
	fmt.Println("Voting Started")
	randomNumber := rand.Intn(100)
	time.Sleep(time.Duration(time.Duration(randomNumber).Microseconds()))
	setRaftState("candidate")

	key := clientID.String() + "term"
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
	candidatePayload = append(candidatePayload, clientID.String(), strconv.Itoa(newTerm))
	token := client.Publish("election/1", 0, false, candidatePayload)
	token.Wait()
}
