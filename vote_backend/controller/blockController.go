package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"vote_backend/models"
	"vote_backend/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CreateBlock() {
	var transactionsCount []models.Transaction
	transactionsTable, err := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	transactionsTable.Find(&transactionsCount)

	fmt.Println("Total transactions in database: " + strconv.Itoa(len(transactionsCount)))
	//todo: if the transactins are less than 5 for a long time, create a block with the few transactions available
	if len(transactionsCount) >= 5 {
		fmt.Println("Creating new block")

		//fetch oldest five rows
		var transactions []models.Transaction

		transactionsTable.Limit(5).Order("created_at asc").Find(&transactions)

		jsonData, err := json.Marshal(transactions)
		if err != nil {
			panic(err)
		}
		fmt.Println("Oldest 5 transactions: " + string(jsonData))

		//create block

		var lastBlock models.Block
		blockTable, err2 := gorm.Open(sqlite.Open("nodeDB.sql"), &gorm.Config{})
		if err2 != nil {
			panic(err2)
		}
		blockTable.Limit(1).Order("created_at desc").Find(&lastBlock)

		type HashableBlock struct {
			Version           int
			PreviousBlockHash string
			Data              string
		}

		hashableBlock := HashableBlock{Version: 1, PreviousBlockHash: lastBlock.BlockHash, Data: string(jsonData)}
		newBlockBytes, err3 := json.Marshal(hashableBlock)
		if err3 != nil {
			panic(err3)
		}
		sum := sha256.Sum256([]byte(newBlockBytes))
		blockHash := hex.EncodeToString(sum[:])

		fmt.Println("Hashable block")
		fmt.Println(hashableBlock)
		newBlock := models.Block{Version: 1, BlockHash: blockHash, PreviousBlockHash: lastBlock.BlockHash, CreatedBy: utils.ReadClientID(), Data: string(jsonData)}
		if err != nil {
			panic(err)
		}
		fmt.Println("New block")
		fmt.Println(newBlock)
		//insert into db
		result := blockTable.Create(&newBlock)
		if result.Error == nil {
			transactionsTable.Delete(&transactions)
		}
		appendToChain(newBlock)

	}

}

func appendToChain(newBlock models.Block) {
	data, err := json.MarshalIndent(newBlock, "", " ")
	if err != nil {
		panic(err)
	}

	chainFile, err := os.OpenFile("chain.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	defer chainFile.Close()
	check_file, err2 := os.Stat("chain.json")
	if err2 != nil {
		panic(err2)
	}
	if check_file.Size() == 0 {
		_, err3 := chainFile.WriteString("[" + "\n" + string(data) + "\n" + "]")
		if err3 != nil {
			panic(err3)
		}
	} else {
		//read the file into an array of structs
		blockChainFile, err4 := os.ReadFile("chain.json")
		if err4 != nil {
			panic(err4)
		}
		var nodesBlocks []models.Block

		err5 := json.Unmarshal(blockChainFile, &nodesBlocks)
		if err5 != nil {
			panic(err5)
		}

		//add the new block to the array
		nodesBlocks = append(nodesBlocks, newBlock)
		data, err6 := json.MarshalIndent(nodesBlocks, "", " ")
		if err6 != nil {
			panic(err6)
		}

		//delete the file
		err7 := os.Remove("chain.json")
		if err7 != nil {
			panic(err7)
		}
		//recreate the file
		newFile, err8 := os.OpenFile("chain.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err8 != nil {
			panic(err8)
		}
		defer newFile.Close()

		//write to the file
		_, err9 := newFile.WriteString("\n" + string(data) + "\n")
		if err9 != nil {
			panic(err9)
		}
	}
}
