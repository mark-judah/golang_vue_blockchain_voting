package controller

import (
	"fmt"
	"net/http"
	"vote_backend/models"

	"github.com/gin-gonic/gin"
)

func NewTransaction(context *gin.Context) {
	fmt.Println("...............New vote", context.Request.Body)

	var newTransaction []models.Transaction
	if err := context.BindJSON(&newTransaction); err != nil {
		return
	}
	for _, x := range newTransaction {
		AppendToLeader(x)
	}
	context.IndentedJSON(http.StatusCreated, newTransaction)
}
