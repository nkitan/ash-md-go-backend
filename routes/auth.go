package routes

import (
    "backend-go/controllers"

    "github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
    authGroup := router.Group("/")
    {
        authGroup.POST("/signup", controllers.Signup)
        authGroup.POST("/signin", controllers.Signin)
        authGroup.POST("/google-auth", controllers.GoogleAuth)
    }
}
