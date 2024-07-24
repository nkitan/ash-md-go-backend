package utils

import (
    "context"
    "log"
    "os"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/joho/godotenv"
)

var S3Client *s3.Client
var PresignClient *s3.PresignClient

func InitAWS() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    cfg, err := config.LoadDefaultConfig(context.TODO(), 
        config.WithRegion(os.Getenv("AWS_REGION")), 
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            os.Getenv("AWS_ACCESS_KEY"), 
            os.Getenv("AWS_SECRET_KEY"), 
            ""),
        ),
    )
    if err != nil {
        log.Fatalf("Unable to load SDK config, %v", err)
    }

    S3Client = s3.NewFromConfig(cfg)
    PresignClient = s3.NewPresignClient(S3Client)
}