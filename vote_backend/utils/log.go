package utils

import (
	"vote_backend/models"
)

type Queue struct {
	Transactions []models.Transaction `json:"transactions"`
}

func (q *Queue) Enqueue(newVote models.Transaction) {
	q.Transactions = append(q.Transactions, newVote)
}
