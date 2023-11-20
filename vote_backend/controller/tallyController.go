package controller

import "github.com/gin-gonic/gin"

func Tally(context *gin.Context) {
	token := Client[0].Publish("tallyVotes/1", 0, false, "tally votes request")
	token.Wait()
}
