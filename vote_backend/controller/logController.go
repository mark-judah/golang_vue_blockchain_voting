package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"vote_backend/models"
)

//var Transactions = make(map[string]models.Transaction)

func Enqueue(newVote models.Transaction) {
	//q.Transactions = append(q.Transactions, newVote)
	print("Enqueue key: " + newVote.Txid)
	// Transactions[newVote.Txid] = newVote

	transactionData, err3 := json.Marshal(newVote)
	if err3 != nil {
		panic(err3)
	}

	token2 := Client[0].Publish("raftLogAppend/1", 0, false, transactionData)
	token2.Wait()
}
func Dequeue(txid string) {
	fmt.Println("Dequeue called")
	//delete(Transactions, txid)

}
func PersistLog(newVote models.Transaction) {
	data, err := json.MarshalIndent(newVote, "", " ")
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile("log.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	check_file, err2 := os.Stat("log.json")
	if err2 != nil {
		panic(err2)
	}
	if check_file.Size() == 0 {
		_, err3 := logFile.WriteString(string(data))
		if err3 != nil {
			panic(err3)
		}
	} else {
		_, err3 := logFile.WriteString("," + string(data))
		if err3 != nil {
			panic(err3)
		}
	}

}
