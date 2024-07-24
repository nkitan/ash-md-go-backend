package main

import (
    "encoding/json"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "backend-go/models"
    "backend-go/routes"
    "backend-go/utils"
)

func main() {
    // Load configuration
    file, err := os.Open("config.json")
    if err != nil {
        log.Fatal("Error opening config file", err)
    }
    defer file.Close()
    var config utils.Config
    
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        log.Fatal("Error decoding config file", err)
    }

    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    // Connect to database
    utils.ConnectDatabase(config)

    // Initialize Firebase
    utils.InitFirebase()

    // Initialize AWS
    utils.InitAWS()

    // Migrate database
    models.UsersMigrate(utils.DB)
    models.BlogsMigrate(utils.DB)
    models.CommentsMigrate(utils.DB)
    models.NotificationsMigrate(utils.DB)

    // Set up Gin router
    router := gin.Default()

    // Define routes
    routes.AuthRoutes(router)
    routes.StorageRoutes(router)
    routes.BlogRoutes(router)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    // Run server
    router.Run(":" + port)
}
