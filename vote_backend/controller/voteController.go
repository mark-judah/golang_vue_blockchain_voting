package controller

import (
	"fmt"
	"net/http"
	"vote_backend/models"

	"github.com/gin-gonic/gin"
)

// the mempool uses a queue datastructure
type Queue struct {
	transactions []models.Vote
}

func NewVote(context *gin.Context) {
	fmt.Println("...............New vote", &context.Request.Body)

	newVote := models.Vote{}
	if err := context.BindJSON(&newVote); err != nil {
		return
	}
	//store vote in transaction pool
	transactionPool := Queue{}

	transactionPool.Enqueue(newVote)
	//inform other nodes
	context.IndentedJSON(http.StatusCreated, newVote)
}

func (q *Queue) Enqueue(newVote models.Vote) {
	q.transactions = append(q.transactions, newVote)
}
