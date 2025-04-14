package controllers

import (
	"demos/internal/services"
	"demos/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) RegisterUser(ctx *gin.Context) {
	var req services.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, response.Error(response.StatusBadRequest, "Invalid request data"))
		return
	}

	user, err := c.userService.RegisterUser(&req)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(response.StatusInternalError, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessWithMessage("User registered successfully", user))
}

func (c *UserController) LoginUser(ctx *gin.Context) {
	var req services.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		ctx.JSON(http.StatusOK, response.Error(response.StatusBadRequest, "Invalid request data"))
		return
	}
	token, err := c.userService.LoginUser(&req)
	if err!= nil {
		ctx.JSON(http.StatusOK, response.Error(response.StatusInternalError, err.Error()))
	}
	ctx.JSON(http.StatusOK, response.SuccessWithMessage("Login successful", token))
}