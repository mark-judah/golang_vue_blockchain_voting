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
var LeaderLogSize = 0

func NewTransaction(context *gin.Context) {
	fmt.Println("...............New vote", context.Request.Body)

	newTransaction := models.Transaction{}
	if err := context.BindJSON(&newTransaction); err != nil {
		return
	}

	Log.Enqueue(newTransaction)
	context.IndentedJSON(http.StatusCreated, newTransaction)
}
