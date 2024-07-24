package controllers

import (
    "context"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/google/uuid"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "backend-go/utils"
)

func GenerateUploadURL(c *gin.Context) {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
    
    imageName := uuid.New().String() + ".jpeg"
    input := &s3.PutObjectInput{
        Bucket: aws.String(os.Getenv("AWS_BUCKET_NAME")),
        Key:    aws.String(imageName),
        ContentType: aws.String("image/jpeg"),
    }

    presignClient := utils.PresignClient
    req, err := presignClient.PresignPutObject(context.TODO(), input, s3.WithPresignExpires(time.Duration(360)*time.Second))
    if err != nil {
        log.Printf("SERVER ERROR [GET-UPLOAD-URL]: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected Error Occurred"})
        return
    }

    log.Printf("Generated URL: %s", req.URL)
    c.JSON(http.StatusOK, gin.H{"uploadURL": req.URL})
}