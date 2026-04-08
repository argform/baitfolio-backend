package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/argform/baitfolio-backend/internal/auth"
	"github.com/argform/baitfolio-backend/internal/service"
	"github.com/argform/baitfolio-backend/internal/transport/http/handlers"
	"github.com/argform/baitfolio-backend/internal/transport/http/middleware"
)

type Dependencies struct {
	AuthService *service.AuthService
	PointService *service.PointService
	JWTManager *auth.JWTManager
}

func NewRouter(deps Dependencies) *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.String(nethttp.StatusOK, "ok")
	})

	authHandler := handlers.NewAuthHandler(deps.AuthService)
	pointHandler := handlers.NewPointHandler(deps.PointService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/me", middleware.Auth(deps.JWTManager), authHandler.GetMe)
	}

	pointGroup := r.Group("/points")
	{
		pointGroup.POST("", middleware.Auth(deps.JWTManager), pointHandler.Create)
		pointGroup.GET("/:id", pointHandler.GetByID)
		pointGroup.GET("", pointHandler.GetAllInsideTile)
	}

	return r
}
