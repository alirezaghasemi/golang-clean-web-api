package routers

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/handlers"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	handler := handlers.NewUsersHandler(cfg)

	router.POST("/send-otp", handler.SendOtp)
}
