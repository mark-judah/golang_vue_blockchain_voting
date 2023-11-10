package controller

import (
	"fmt"
	"net/http"
	"vote_backend/models"
	"vote_backend/utils"

	"github.com/gin-gonic/gin"
)

// the log uses a queue datastructure
var Log = utils.Queue{}

func NewVote(context *gin.Context) {
	fmt.Println("...............New vote", context.Request.Body)

	newVote := models.Vote{}
	if err := context.BindJSON(&newVote); err != nil {
		return
	}
	//leader node sends an append request to the follower nodes logs
	//follower nodes update logs with the clients request(newVote)
	//leader node receives confirmation from majority of the nodes
	//leader node updates its log with the newVote
	Log.Enqueue(newVote)

	context.IndentedJSON(http.StatusCreated, newVote)
}
