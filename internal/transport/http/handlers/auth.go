package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/argform/baitfolio-backend/internal/domain"
	"github.com/argform/baitfolio-backend/internal/service"
	httpresponse "github.com/argform/baitfolio-backend/internal/transport/http/response"
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
	Username string `json:"username"`
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
	Username string `json:"username"`
	Email string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName *string `json:"last_name"`
	About *string `json:"about"`
}

func newUserResponse(user *domain.User) UserResponse {
	safeUser := user.Sanitized()

	return UserResponse{
		UserID: safeUser.UserID,
		Username: safeUser.Username,
		Email: safeUser.Email,
		FirstName: safeUser.FirstName,
		LastName: safeUser.LastName,
		About: safeUser.About,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid request body")
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
		httpresponse.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusCreated, "ok")
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Login(c.Request.Context(), service.LoginInput{
		Email: req.Email,
		Password: req.Password,
	})
	if err != nil {
		httpresponse.WriteError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		Token: token,
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		httpresponse.WriteError(c, http.StatusBadRequest, "missing user context")
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		httpresponse.WriteError(c, http.StatusBadRequest, "invalid user context")
		return
	}

	user, err := h.authService.GetMe(c.Request.Context(), userID)
	if err != nil {
		httpresponse.WriteError(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, newUserResponse(user))
}
