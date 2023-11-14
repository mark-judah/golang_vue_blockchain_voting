package controller

import (
	"fmt"
	"net/http"
	"vote_backend/models"

	"github.com/gin-gonic/gin"
)

// the log uses a queue datastructure for temporary storage and a json file for permanent storage

func NewTransaction(context *gin.Context) {
	fmt.Println("...............New vote", context.Request.Body)

	newTransaction := models.Transaction{}
	if err := context.BindJSON(&newTransaction); err != nil {
		return
	}

	Enqueue(newTransaction)

	context.IndentedJSON(http.StatusCreated, newTransaction)
}
