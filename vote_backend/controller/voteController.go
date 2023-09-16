package controller

import (
	"fmt"
	"net/http"
	"vote_backend/models"
	"vote_backend/utils"

	"github.com/gin-gonic/gin"
)

// the mempool uses a queue datastructure
var TransactionPool = utils.Queue{}

func NewVote(context *gin.Context) {
	fmt.Println("...............New vote", context.Request.Body)

	newVote := models.Vote{}
	if err := context.BindJSON(&newVote); err != nil {
		return
	}
	//store vote in transaction pool
	TransactionPool.Enqueue(newVote)

	context.IndentedJSON(http.StatusCreated, newVote)
}
