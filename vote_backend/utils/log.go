package utils

import (
	"vote_backend/models"
)

type Queue struct {
	Transactions []models.Vote `json:"transactions"`
}

func (q *Queue) Enqueue(newVote models.Vote) {
	q.Transactions = append(q.Transactions, newVote)
}
