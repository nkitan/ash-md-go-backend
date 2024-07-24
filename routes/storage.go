package routes

import (
    "backend-go/controllers"
    "backend-go/middleware"

    "github.com/gin-gonic/gin"
)

func StorageRoutes(router *gin.Engine) {
    storageGroup := router.Group("/")
    storageGroup.Use(middleware.Authenticate())
    {
        storageGroup.GET("/get-upload-url", controllers.GenerateUploadURL)
    }
}