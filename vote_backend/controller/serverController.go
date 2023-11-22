package controller

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartApiServer() {
	fmt.Println("Starting API server")
	//only the leader can create a router and receive requests
	//if the server is unreachable, the leader is probably dead
	router := gin.Default()

	//setup cors
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	api := router.Group("/api")
	{
		api.POST("/new-vote", NewTransaction)
		api.POST("/tally-votes", Tally)
		api.POST("/create-user", CreateUser)
		api.POST("/login", Login)

		securedRoutes := api.Group("/secured").Use(Auth())
		{
			securedRoutes.GET("/current-user", CurrentUser)
			securedRoutes.GET("/get-all-users", GetUsers)

		}
	}
	router.Run("localhost:3500")
}
