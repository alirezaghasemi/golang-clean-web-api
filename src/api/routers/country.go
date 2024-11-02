package routers

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/handlers"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/gin-gonic/gin"
)

func Country(r *gin.RouterGroup, cfg *config.Config) {
	handler := handlers.NewCountryHandler(cfg)

	r.POST("/", handler.Create)
	r.PUT("/:id", handler.Update)
	r.DELETE("/:id", handler.Delete)
	r.GET("/:id", handler.GetById)
}
