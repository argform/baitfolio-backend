package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/argform/baitfolio-backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Username  string  `json:"username"`
	Email string  `json:"email"`
	Password string  `json:"password"`
	FirstName *string `json:"first_name"`
	LastName *string `json:"last_name"`
	About *string `json:"about"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type UserResponse struct {
	UserID uint64 `json:"user_id"`
	Username string  `json:"username"`
	Email string  `json:"email"`
	FirstName *string `json:"first_name"`
	LastName *string `json:"last_name"`
	About *string `json:"about"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	_, err := h.authService.Register(c.Request.Context(), service.RegisterInput{
		Username: req.Username,
		Email: req.Email,
		Password: req.Password,
		FirstName: req.FirstName,
		LastName: req.LastName,
		About: req.About,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, "ok")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), service.LoginInput{
		Email: req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		Token: token,
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing user context",
		})
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user context",
		})
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		UserID: user.UserID,
		Username: user.Username,
		Email: user.Email,
		FirstName: user.FirstName,
		LastName: user.LastName,
		About: user.About,
	})
}
