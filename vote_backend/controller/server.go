package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func StartApiServer() {
	fmt.Println("Starting API server")
	//only the leader can create a router and receive requests
	//if the server is unreachable, the leader is probably dead
	router := gin.Default()
	router.POST("/new-vote", NewTransaction)
	router.Run("localhost:3500")
}
