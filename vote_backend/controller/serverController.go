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
	config.MaxAge = 1 * time.Hour

	router.Use(cors.New(config))
	api := router.Group("/api")
	{
		api.POST("/tally-votes", Tally)
		api.POST("/login", Login)

		securedRoutes := api.Group("/secured").Use(Auth())
		{
			securedRoutes.POST("/create-user", CreateUser)
			securedRoutes.GET("/current-user", CurrentUser)
			securedRoutes.GET("/get-all-users", GetUsers)
			securedRoutes.POST("/new-county", NewCounty)
			securedRoutes.GET("/get-all-counties", FetchCounties)
			securedRoutes.POST("/new-constituency", NewConstituency)
			securedRoutes.GET("/get-all-constituencies", FetchConstituencies)
			securedRoutes.POST("/new-ward", NewWard)
			securedRoutes.GET("/get-all-wards", FetchWards)
			securedRoutes.POST("/new-polling-station", NewPollingStation)
			securedRoutes.GET("/get-all-polling-stations", FetchPollingStations)
			securedRoutes.POST("/new-desktop-client", NewDesktopClient)
			securedRoutes.GET("/get-all-desktop-clients", FetchDesktopClients)
			securedRoutes.POST("/new-vote", NewTransaction)
			securedRoutes.GET("/get-all-candidates", FetchCandidates)
			securedRoutes.POST("/new-candidate", NewCandidate)
			securedRoutes.GET("/get-all-voters", FetchVoters)
			securedRoutes.POST("/new-voter", NewVoter)

		}
	}
	router.Run("localhost:3500")
}
