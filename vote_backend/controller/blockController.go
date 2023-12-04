package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"unsafe"
	"vote_backend/models"
	"vote_backend/utils"

	"github.com/gin-gonic/gin"
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

func FetchBlockChain(context *gin.Context) {
	blockChainFile, err := os.ReadFile("chain.json")
	if err != nil {
		fmt.Println(err)
	}
	type blockData struct {
		BlockHeight       int
		BlockHash         string
		PreviousBlockHash string
		NoOfTransactions  int
		BlockSize         uintptr
	}
	var nodesBlocks []models.Block
	var data []blockData

	err2 := json.Unmarshal(blockChainFile, &nodesBlocks)
	if err2 != nil {
		fmt.Println(err2)
	}

	for _, x := range nodesBlocks {
		y := x.Data
		var transactions []models.Transaction
		err3 := json.Unmarshal([]byte(y), &transactions)
		if err3 != nil {
			panic(err3)
		}
		data = append(data, blockData{BlockHeight: x.Index, BlockHash: x.BlockHash, PreviousBlockHash: x.PreviousBlockHash, NoOfTransactions: len(transactions), BlockSize: unsafe.Sizeof(x)})
	}
	context.IndentedJSON(http.StatusOK, data)

}

func FetchNetworkState(context *gin.Context) {
	blockChainFile, err := os.ReadFile("chain.json")
	if err != nil {
		fmt.Println(err)
	}
	type blockData struct {
		BlockHeight      int
		NoOfTransactions int
	}
	var nodesBlocks []models.Block

	err2 := json.Unmarshal(blockChainFile, &nodesBlocks)
	if err2 != nil {
		fmt.Println(err2)
	}

	var transactions []models.Transaction
	var all_transactions []models.Transaction

	for _, x := range nodesBlocks {
		y := x.Data
		err3 := json.Unmarshal([]byte(y), &transactions)
		if err3 != nil {
			panic(err3)
		}
		all_transactions = append(all_transactions, transactions...)
	}

	data := blockData{BlockHeight: len(nodesBlocks), NoOfTransactions: len(all_transactions)}
	context.IndentedJSON(http.StatusOK, data)

}

func FindTransactionByID(context *gin.Context) {
	type txid struct {
		txid string
	}
	var txidVal txid
	if err := context.BindJSON(&txidVal); err != nil {
		return
	}

	blockChainFile, err := os.ReadFile("chain.json")
	if err != nil {
		fmt.Println(err)
	}
	type transactionData struct {
		BlockHeight int
		Transaction models.Transaction
	}

	var data []transactionData
	var nodesBlocks []models.Block

	err2 := json.Unmarshal(blockChainFile, &nodesBlocks)
	if err2 != nil {
		fmt.Println(err2)
	}

	var transaction models.Transaction

	for i, x := range nodesBlocks {
		y := x.Data
		err3 := json.Unmarshal([]byte(y), &transaction)
		if err3 != nil {
			panic(err3)
		}
		if transaction.Txid == txidVal.txid {
			data = append(data, transactionData{BlockHeight: i, Transaction: transaction})
		}
	}
	context.IndentedJSON(http.StatusOK, data)

}

func FindTransactionsByBlockHash(context *gin.Context) {

	type hash struct {
		hash string
	}
	var hashVal hash
	if err := context.BindJSON(&hashVal); err != nil {
		return
	}

	blockChainFile, err := os.ReadFile("chain.json")
	if err != nil {
		fmt.Println(err)
	}

	var nodesBlocks []models.Block

	err2 := json.Unmarshal(blockChainFile, &nodesBlocks)
	if err2 != nil {
		fmt.Println(err2)
	}

	var transactions []models.Transaction

	for _, x := range nodesBlocks {

		if x.BlockHash == hashVal.hash {
			y := x.Data
			err3 := json.Unmarshal([]byte(y), &transactions)
			if err3 != nil {
				panic(err3)
			}
		}

	}
	context.IndentedJSON(http.StatusOK, transactions)

}
