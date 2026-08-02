package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	authdto "github.com/murashi19/koda-b8-backend/internal/dto/auth"
	"github.com/murashi19/koda-b8-backend/internal/lib/response"
	"github.com/murashi19/koda-b8-backend/internal/service"
)

type AuthHandler struct {
	service *service.UserService
}

func NewAuthHandler(service *service.UserService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req authdto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err := h.service.Register(c.Request.Context(), &req)
	fmt.Println("error disini", err)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Register success", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authdto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Login success", result)
}
