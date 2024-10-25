package api

import (
	"fmt"
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/middlewares"
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/routers"
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/validations"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func InitServer() {

	cfg := config.GetConfig()

	r := gin.New()

	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		val.RegisterValidation("mobile", validations.IranianMobileNumberValidator, true)
		val.RegisterValidation("password", validations.PasswordValidator, false)
	}

	r.Use(middlewares.Cors(cfg))
	r.Use(gin.Logger(), gin.Recovery(), middlewares.LimitByRequest())

	api := r.Group("/api")

	v1 := api.Group("/v1")
	{
		health := v1.Group("/health")
		routers.Health(health)

	}

	{
		testHealth := v1.Group("/test")
		routers.TestHealth(testHealth)
	}

	//r.Run(":5005")
	r.Run(fmt.Sprintf(":%s", cfg.Server.Port))
}
