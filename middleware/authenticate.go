package middleware

import (
    "net/http"
    "os"
    "strings"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/golang-jwt/jwt/v4"
)

func Authenticate() gin.HandlerFunc {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User token not found"})
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "User token not found"})
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
                c.Abort()
                return nil, nil
            }
            return []byte(os.Getenv("ACCESS_TOKEN_SECRET")), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            if user, ok := claims["id"].(string); ok {
                c.Set("user", user)
            } else {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
                c.Abort()
                return
            }
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        c.Next()
    }
}
