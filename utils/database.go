package utils

import (
    "fmt"
    "log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string
}

func ConnectDatabase(config Config) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
        config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database: ", err)
    }
    DB = database
    log.Println("Database connection established")
}
