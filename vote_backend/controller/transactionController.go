package controller

import (
	"fmt"
	"net/http"
	"vote_backend/models"

	"github.com/gin-gonic/gin"
)

func NewTransaction(context *gin.Context) {
	fmt.Println("...............New vote", context.Request.Body)

	newTransaction := models.Transaction{}
	if err := context.BindJSON(&newTransaction); err != nil {
		return
	}

	AppendToLeader(newTransaction)

	context.IndentedJSON(http.StatusCreated, newTransaction)
}
