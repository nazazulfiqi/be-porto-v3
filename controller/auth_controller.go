package controller

import (
	"be-porto-v3/dto"
	"be-porto-v3/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController menangani request HTTP untuk autentikasi
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController membuat instance baru AuthController
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService}
}

// Register menangani endpoint untuk registrasi user
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := c.authService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := c.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
