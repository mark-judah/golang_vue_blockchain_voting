package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Authorization")
		cleanToken := strings.TrimSpace(strings.ReplaceAll(token, "Bearer", ""))
		if cleanToken == "" {
			context.IndentedJSON(http.StatusUnauthorized, "Unauthorized")
			context.Abort()
			return
		}
		err := ValidateToken(cleanToken)
		if err != nil {
			context.IndentedJSON(http.StatusUnauthorized, err.Error())
			context.Abort()
			return
		}
	}
}
