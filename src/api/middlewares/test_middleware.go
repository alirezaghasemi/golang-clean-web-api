package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("x-api-key")
		if apiKey == "1" {
			ctx.Next()
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"result": "Api key is required",
		})

		return
	}
}
