package controller

import (
	"vote_backend/models"

	"github.com/gin-gonic/gin"
)

// use postgres to store the data

// routes
// new voter
// new candidate
// new voting node
// new election official
// fetchTransactionPool
// verify that all requests are coming from an official node
func NewVoter(context *gin.Context) {
	newVoter := models.Voter{}
	if err := context.BindJSON(&newVoter); err != nil {
		return
	}

}

func NewVotingNode(context *gin.Context) {
}

func NewElectionOfficial(context *gin.Context) {
}

func NewCandidate(context *gin.Context) {
}

func FetchTransactionPool(context *gin.Context) {
}
