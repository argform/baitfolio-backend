package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/argform/baitfolio-backend/internal/auth"
	"github.com/argform/baitfolio-backend/internal/service"
	"github.com/argform/baitfolio-backend/internal/transport/http/handlers"
)

type Dependencies struct {
	AuthService *service.AuthService
	JWTManager  *auth.JWTManager
}

func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.String(nethttp.StatusOK, "ok")
	})

	authHandler := handlers.NewAuthHandler(deps.AuthService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	return r
}