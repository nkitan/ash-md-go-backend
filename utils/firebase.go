package utils

import (
    "context"
    "log"
    "os"

    firebase "firebase.google.com/go/v4"
    "firebase.google.com/go/v4/auth"
    "google.golang.org/api/option"
    "github.com/joho/godotenv"
)

var FirebaseAuth *auth.Client

func InitFirebase() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
    //opt := option.WithCredentialsFile("/home/nkitan/go/src/backend-go/credentials/storageAccountService.json")
    opt := option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS"))

    app, err := firebase.NewApp(context.Background(), nil, opt)
    if err != nil {
        log.Fatalf("Error initializing app: %v\n", err)
    }

    FirebaseAuth, err = app.Auth(context.Background())
    if err != nil {
        log.Fatalf("Error getting Auth client: %v\n", err)
    }
}
