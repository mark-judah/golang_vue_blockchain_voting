package utils

import (
	"encoding/json"
	"os"
	"vote_backend/models"
)

type Queue struct {
	Transactions []models.Transaction `json:"transactions"`
}

func (q *Queue) Enqueue(newVote models.Transaction) {
	q.Transactions = append(q.Transactions, newVote)
	PersistLog(newVote)
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
