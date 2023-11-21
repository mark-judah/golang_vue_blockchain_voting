package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"vote_backend/models"

	"github.com/gin-gonic/gin"
)

func Tally(context *gin.Context) {

	//change to self excecuting function when a new block is added to the chain
	token := Client[0].Publish("tallyVotes/1", 0, false, "pass context here")
	token.Wait()
}

func TallyFeedback() {
	blockChainFile, err := os.ReadFile("chain.json")
	if err != nil {
		panic(err)
	}
	var nodesBlocks []models.Block

	err2 := json.Unmarshal(blockChainFile, &nodesBlocks)
	if err2 != nil {
		panic(err2)
	}

	type HashableBlock struct {
		Version           int
		PreviousBlockHash string
		Data              string
	}
	var hashableBlock HashableBlock
	//verify block data
	for i, x := range nodesBlocks {
		hashableBlock.Version = x.Version
		hashableBlock.PreviousBlockHash = x.PreviousBlockHash
		hashableBlock.Data = x.Data

		fmt.Println(hashableBlock)
		blockBytes, err := json.Marshal(hashableBlock)
		if err != nil {
			panic(err)
		}
		sum := sha256.Sum256([]byte(blockBytes))
		blockHash := hex.EncodeToString(sum[:])
		fmt.Println("Block " + strconv.Itoa(i) + " calculated hash: " + blockHash)
		fmt.Println("vs")
		fmt.Println("Block " + strconv.Itoa(i) + " stored hash: " + x.BlockHash)

		if blockHash == x.BlockHash {
			fmt.Println("Block " + strconv.Itoa(i) + " is valid")
		} else {
			fmt.Println("Block " + strconv.Itoa(i) + " is invalid")
			//delete the log.json file
			//delete the database file
			//delete the chain.json file
			//set nodeSyncCounter to zero, the node will resync to generate a valid chain from the network
			break
		}
	}

	//verify chain order
	hashes := []string{}
	for i, x := range nodesBlocks {
		if x.PreviousBlockHash == "" {
			fmt.Println("First block")
			hashes = append(hashes, x.BlockHash)
			continue
		}
		if x.PreviousBlockHash == hashes[i-1] {
			hashes = append(hashes, x.BlockHash)
			fmt.Println("Block " + strconv.Itoa(i) + " order is valid")

		} else {
			fmt.Println("Block " + strconv.Itoa(i) + " is invalid")
			//delete the log.json file
			//delete the database file
			//delete the chain.json file
			//set nodeSyncCounter to zero, the node will resync to generate a valid chain from the network
			break
		}
	}

	//count results
	var results []string
	var transaction []models.Transaction

	for _, x := range nodesBlocks {
		data := x.Data

		err2 := json.Unmarshal([]byte(data), &transaction)
		if err2 != nil {
			panic(err2)
		}
		for _, y := range transaction {
			results = append(results, y.CandidateId)
		}
	}
	var tally = make(map[string]int)
	for _, candidate := range results {
		tally[candidate]++
	}

	fmt.Println("results")
	fmt.Println(tally)

}
