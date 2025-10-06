package routes

import (
    "row-level-security-backend/handlers"
    "row-level-security-backend/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
    authHandler := handlers.NewAuthHandler(db)
    postHandler := handlers.NewPostHandler(db)

    api := router.Group("/api/v1")
    {
        auth := api.Group("/auth")
        {
            auth.POST("/login", authHandler.Login)
            auth.POST("/register", authHandler.Register)
        }
    }

    protected := api.Group("/")
    protected.Use(middleware.AuthMiddleware())
    protected.Use(middleware.RLSMiddleware())
    {
        posts := protected.Group("/posts")
        {
            posts.POST("", postHandler.CreatePost)
            posts.GET("", postHandler.GetPosts)
            posts.GET("/:id", postHandler.GetPost)
            posts.PUT("/:id", postHandler.UpdatePost)
            posts.DELETE("/:id", postHandler.DeletePost)
        }
    }
}