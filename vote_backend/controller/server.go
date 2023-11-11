package controller

import (
	"github.com/gin-gonic/gin"
)

func StartApiServer() {
	//only the leader can create a router and receive requests
	//if the server is unreachable, the leader is probably dead
	router := gin.Default()
	router.POST("/new-vote", NewTransaction)
	router.Run("localhost:3500")
}
