package handlers

import (
	"fmt"
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck godoc
// @Summary Health Check
// @Description Health Check
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/health/ [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, helper.GenerateBaseResponse("Working!", true, 0))
	return
}

func (h *HealthHandler) HealthPost(c *gin.Context) {
	c.JSON(http.StatusOK, "Working Post!")
	return
}

func (h *HealthHandler) HealthPostById(c *gin.Context) {
	//id := c.Param("id")
	id := c.Params.ByName("id")
	c.JSON(http.StatusOK, fmt.Sprintf("Working Post By ID: %s", id))
	return
}
