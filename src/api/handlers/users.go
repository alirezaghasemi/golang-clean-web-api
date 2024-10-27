package handlers

import (
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/dto"
	"github.com/alirezaghasemi/golang-clean-web-api/src/api/helper"
	"github.com/alirezaghasemi/golang-clean-web-api/src/config"
	"github.com/alirezaghasemi/golang-clean-web-api/src/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	service *services.UserService
}

func NewUsersHandler(cfg *config.Config) *UserHandler {
	service := services.NewUserService(cfg)
	return &UserHandler{service: service}
}

// SendOtp godoc
// @Summary Send otp to user
// @Description Send otp to user
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.GetOtpRequest true "get otp request"
// @Success 201 {object} helper.BaseHttpResponse "Success"
// @Failure 400 {object} helper.BaseHttpResponse "Failed"
// @Failure 409 {object} helper.BaseHttpResponse "Failed"
// @Router /v1/users/send-otp [post]
func (u *UserHandler) SendOtp(ctx *gin.Context) {
	req := new(dto.GetOtpRequest)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, helper.GenerateBaseResponseWithValidationError(nil, false, -1, err))
		return
	}

	err = u.service.SendOtp(req)
	if err != nil {
		ctx.AbortWithStatusJSON(helper.TranslateErrorToStatusCode(err), helper.GenerateBaseResponseWithError(nil, false, -1, err))
		return
	}

	// Call internal sms service
	ctx.JSON(http.StatusCreated, helper.GenerateBaseResponse(nil, true, 0))
}
