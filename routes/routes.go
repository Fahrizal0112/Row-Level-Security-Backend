package routes

import (
	"row-level-security-backend/handlers"
	"row-level-security-backend/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	authHandler := handlers.NewAuthHandler(db)
	postsHandler := handlers.NewPostHandler(db)
	tenantHandler := handlers.NewTenantHandler(db)

	api := r.Group("/api/v1")
	
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	tenant := api.Group("/tenant")
	tenant.Use(middleware.AuthMiddleware())
	{
		tenant.POST("/", tenantHandler.CreateTenant)
		tenant.GET("/", tenantHandler.GetMyTenant)
	}

	posts := api.Group("/posts")
	posts.Use(middleware.AuthMiddleware())
	{
		posts.POST("/", postsHandler.CreatePost)
		posts.GET("/", postsHandler.GetPosts)
		posts.GET("/:id", postsHandler.GetPost)
		posts.PUT("/:id", postsHandler.UpdatePost)
		posts.DELETE("/:id", postsHandler.DeletePost)
	}
}
