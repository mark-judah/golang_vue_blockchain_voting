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

func (q *Queue) Dequeue() models.Vote {
	//call after a vote is verified and added to the blockchain
	dequeued := q.Transactions[0]
	//set the queue to include 1 to the last index and exclude the index 0
	q.Transactions = q.Transactions[1:]
	return dequeued
}
