package models

import (
	"time"
)

type Transaction struct {
	Txid            string    `json:"txid" gorm:"primaryKey"`
	NodeId          string    `json:"nodeId"`
	CandidateId     string    `json:"candidateId"`
	CreatedAt       time.Time `json:"timestamp"`
	TransactionHash string    `json:"transactionHash"`
}
