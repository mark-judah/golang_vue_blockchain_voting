package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
	"vote_backend/controller"
	"vote_backend/models"
	"vote_backend/utils"
)

func main() {

	controller.InitMqttClient()
	controller.InitSqlite()
	utils.SetRaftState("follower")
	utils.SetRaftTerm(0)
	utils.SetVoteAndTerm("0", "0", "0")
	//for the admin panel(vue)
	go startHttpServer()

	var nodeSyncCounter = 0
	for {

		if nodeSyncCounter == 1 {
			controller.NodeSync()
		}
		key := utils.ReadClientID() + "state"
		val, err := utils.RedisClient.Get(utils.Ctx, key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("\n My state----->" + val)
		if !controller.LeaderAlive {
			controller.LeaderAliveCounter = controller.LeaderAliveCounter + 1
		}
		if controller.LeaderAliveCounter >= 10 && val == "follower" {
			fmt.Println("Leader Dead")
			requestVotes()

		}

		if val == "leader" {
			key := utils.ReadClientID() + "term"
			val, err := utils.RedisClient.Get(utils.Ctx, key).Result()
			if err != nil {
				panic(err)
			}

			myPayload := append(controller.LeaderPayload, val, "Leader Alive")

			mqttMessage := models.Message{Type: "leader_term_pulse", Payload: myPayload}
			data, err := json.Marshal(mqttMessage)
			if err != nil {
				panic(err)
			}

			time.Sleep(time.Duration(time.Second))
			// publish a message every one second
			token := controller.Client[0].Publish("leaderNodePulse/"+key, 0, false, data)
			token.Wait()

		}

		if val == "candidate" {
			//set to zero to avoid requestVotes being called infinitely
			controller.LeaderAliveCounter = 0
			go killApiServer()

		}

		if val == "follower" {
			fmt.Println("\n Leader Alive--------------------->" + strconv.FormatBool(controller.LeaderAlive))
			fmt.Println("\n Leader Alive Counter--------------------->" + strconv.Itoa(controller.LeaderAliveCounter))
			fmt.Println(utils.GetClientState())
			go killApiServer()

		}

		time.Sleep(time.Duration(time.Second))

		fmt.Println("\n My id----->" + utils.ReadClientID())

		if controller.LeaderAlive {
			nodeSyncCounter = nodeSyncCounter + 1
		}

		go controller.CreateBlock()

	}

}

func requestVotes() {
	fmt.Println("\n --------------------->" + "Voting Started")
	randomNumber := rand.Intn(10)
	fmt.Println("\n --------------------->" + "Waited for" + strconv.Itoa(randomNumber))
	time.Sleep(time.Duration(time.Duration(randomNumber).Seconds()))
	utils.SetRaftState("candidate")

	key := utils.ReadClientID() + "term"
	term, err := utils.RedisClient.Get(utils.Ctx, key).Result()
	if err != nil {
		panic(err)
	}

	newTerm, err2 := strconv.Atoi(term)
	if err != nil {
		panic(err2)
	}
	newTerm = newTerm + 1
	utils.SetRaftTerm(newTerm)

	var candidatePayload []string
	candidatePayload = append(candidatePayload, utils.ReadClientID(), strconv.Itoa(newTerm))

	mqttMessage := models.Message{Type: "node_election", Payload: candidatePayload}
	data, err := json.Marshal(mqttMessage)
	if err != nil {
		panic(err)
	}
	token := controller.Client[0].Publish("election/1", 0, false, data)
	token.Wait()
}

func startHttpServer() {
	for port := 8080; port <= 9080; port++ {
		ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

		if err != nil {
			fmt.Println("Port in use:", err)
		} else {
			_ = ln.Close()
			fmt.Println("Open port at: " + strconv.Itoa(port))

			cmd := exec.Command("npm", "run", "serve", "--", "--port", strconv.Itoa(port))
			cmd.Dir = "ui"
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			if err != nil {
				fmt.Printf("%s", err)
			} else {
				utils.SetHttpPort(strconv.Itoa(port))
			}
			break
		}
	}

}

func killApiServer() {
	//kill the api server so that only the leader node receives api requests
	command := "fuser -n tcp -k 3500"

	out, err := exec.Command(command).Output()
	if err != nil {
		fmt.Printf("%s", err)
		fmt.Println("Command Successfully Executed")
		output := string(out[:])
		fmt.Println(output)
		print("Api server killed")
	}
}
