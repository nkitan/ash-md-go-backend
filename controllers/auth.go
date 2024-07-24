package controllers

import (
    "backend-go/models"
    "backend-go/utils"
    "context"
    "net/http"
    "os"
    "strings"
    "time"
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/google/uuid"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

func formatDataToSend(user models.User) map[string]interface{} {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":  user.ID,
        "exp": time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
    if err != nil {
        panic(err) // Handle the error appropriately in production
    }

    return map[string]interface{}{
        "access_token": tokenString,
        "profile_img":  user.PersonalInfo.ProfileImg,
        "username":     user.PersonalInfo.Username,
        "fullname":     user.PersonalInfo.Fullname,
    }
}

func generateUsername(email string, db *gorm.DB) (string, error) {
    username := strings.Split(email, "@")[0]
    var count int64
    db.Model(&models.User{}).Where("username = ?", username).Count(&count)
    if count > 0 {
        username += uuid.NewString()[:5]
    }
    return username, nil
}

func Signup(c *gin.Context) {
    var input struct {
        Fullname string `json:"fullname" binding:"required,min=3"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
        return
    }

    username, err := generateUsername(input.Email, utils.DB)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
        return
    }

    user := models.User{
        PersonalInfo: models.PersonalInfo{
            Fullname: input.Fullname,
            Email:    input.Email,
            Password: string(hashedPassword),
            Username: username,
            ProfileImg: utils.GetDefaultProfileImg(),
        },
        SocialLinks: models.SocialLinks{
            Youtube:   "",
            Instagram: "",
            Facebook:  "",
            Twitter:   "",
            Github:    "",
            Website:   "",
        },
        Blogs: []models.Blog{}, // Initialize as an empty slice of UUIDs
    }

    if err := utils.DB.Create(&user).Error; err != nil {
        if strings.Contains(err.Error(), "duplicate key value") {
            c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
        return
    }

    c.JSON(http.StatusOK, formatDataToSend(user))
}

func Signin(c *gin.Context) {
    var input struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := utils.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    if user.GoogleAuth {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "This account was created using Google. Please log in using Google"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PersonalInfo.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    c.JSON(http.StatusOK, formatDataToSend(user))
}

func GoogleAuth(c *gin.Context) {
    var input struct {
        AccessToken string `json:"access_token" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, err := utils.FirebaseAuth.VerifyIDToken(context.Background(), input.AccessToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    email := token.Claims["email"].(string)
    name := token.Claims["name"].(string)
    picture := strings.Replace(token.Claims["picture"].(string), "s96-c", "s384-c", 1)

    var user models.User
    if err := utils.DB.Where("email = ?", email).First(&user).Error; err != nil {
        if !strings.Contains(err.Error(), "record not found") {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
            return
        }

        username, err := generateUsername(email, utils.DB)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
            return
        }

        user = models.User{
            PersonalInfo: models.PersonalInfo{
                Fullname:   name,
                Email:      email,
                ProfileImg: picture,
                Username:   username,
            },
            GoogleAuth: true,
        }

        if err := utils.DB.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
            return
        }
    } else if !user.GoogleAuth {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "This account was created without Google authentication. Please log in with a password to access the account"})
        return
    }

    c.JSON(http.StatusOK, formatDataToSend(user))
}
