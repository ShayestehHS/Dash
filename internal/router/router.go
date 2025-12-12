package router

import (
	"github.com/gin-gonic/gin"

	userCore "Dash/core/user"
	userRepo "Dash/db/repository/user/postgresql"
	userAPI "Dash/internal/api/user"
	"Dash/internal/database"
	"Dash/pkg/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"error": "method not allowed"})
	})

	r.Use(middleware.ErrorHandler())
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health/", func(c *gin.Context) {
			if err := database.HealthCheck(); err != nil {
				c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "healthy"})
		})

		setupUserRoutes(api)
	}

	return r
}

func setupUserRoutes(api *gin.RouterGroup) {
	repo := userRepo.NewRepository(database.DB)
	service := userCore.NewService(repo)
	handler := userAPI.NewHandler(service)

	user := api.Group("/user")
	{
		auth := user.Group("/auth")
		{
			auth.POST("/login/", handler.Login)
		}
	}
}
