package routers

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/handlers"
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/middlewares"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/gin-gonic/gin"
)

func User(router *gin.RouterGroup, cfg *config.Config) {
	handler := handlers.NewUsersHandler(cfg)

	router.POST("/send-otp", middlewares.OtpLimiter(cfg), handler.SendOtp)
	router.POST("/login-by-username", handler.LoginByUsername)
	router.POST("/login-by-mobile", handler.RegisterLoginByMobileNumber)
	router.POST("/register-by-username", handler.RegisterByUsername)
}
