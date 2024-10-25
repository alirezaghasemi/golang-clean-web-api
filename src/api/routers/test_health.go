package routers

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/handlers"
	"github.com/gin-gonic/gin"
)

func TestHealth(r *gin.RouterGroup) {
	handler := handlers.NewTestHandler()

	r.GET("/", handler.Test)
	r.POST("/binder/header1", handler.HeaderBinder1)
	r.POST("/binder/header2", handler.HeaderBinder2)

	r.POST("/binder/query1", handler.QueryBinder1)
	r.POST("/binder/query2", handler.QueryBinder2)

	r.POST("/binder/uri/:id/:name", handler.UriBinder)

	r.POST("/binder/body", handler.BodyBinder)

	r.POST("/binder/form", handler.FormBinder)

	r.POST("/binder/file", handler.FileBinder)
}
