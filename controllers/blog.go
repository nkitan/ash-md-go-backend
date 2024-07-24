package controllers

import (
    "backend-go/models"
    "backend-go/utils"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "fmt"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    "github.com/google/uuid"
    "github.com/matoous/go-nanoid/v2"
)

const (
    latestBlogsLimit   = 5
    trendingBlogsLimit = 5
    searchBlogsLimit   = 5
)

func GetLatestBlogs(c *gin.Context) {
    page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
    if err != nil || page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }

    var blogs []models.Blog
    if err := utils.DB.
        Where("draft = ?", false).
        Preload("Author", func(db *gorm.DB) *gorm.DB {
            return db.Select("personal_info.profile_img, personal_info.username, personal_info.fullname")
        }).
        Order("published_at desc").
        Select("blog_id, title, des, banner, activity, tags, published_at").
        Offset((page - 1) * latestBlogsLimit).
        Limit(latestBlogsLimit).
        Find(&blogs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch latest blogs"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

func GetLatestBlogsCount(c *gin.Context) {
    var count int64
    if err := utils.DB.
        Model(&models.Blog{}).
        Where("draft = ?", false).
        Count(&count).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch latest blog count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"totalDocs": count})
}

func GetTrendingBlogs(c *gin.Context) {
    var blogs []models.Blog
    if err := utils.DB.
        Where("draft = ?", false).
        Preload("Author", func(db *gorm.DB) *gorm.DB {
            return db.Select("personal_info.profile_img, personal_info.username, personal_info.fullname")
        }).
        Order("activity.total_reads desc, activity.total_likes desc, published_at desc").
        Select("blog_id, title, published_at").
        Limit(trendingBlogsLimit).
        Find(&blogs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch trending blogs"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

func GetSearchBlogsCount(c *gin.Context) {
    var input struct {
        Tag   string `json:"tag"`
        Query string `json:"query"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var count int64
    db := utils.DB.Model(&models.Blog{}).Where("draft = ?", false)

    if input.Tag != "" {
        db = db.Where("tags @> ?", "{"+input.Tag+"}")
    } else if input.Query != "" {
        db = db.Where("title ILIKE ?", "%"+input.Query+"%")
    }

    if err := db.Count(&count).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch search blog count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"totalDocs": count})
}

func SearchBlogs(c *gin.Context) {
    var input struct {
        Tag   string `json:"tag"`
        Query string `json:"query"`
        Page  int    `json:"page" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.Page < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
        return
    }

    var blogs []models.Blog
    db := utils.DB.Where("draft = ?", false)

    if input.Tag != "" {
        // Convert tag to JSONB array
        tagFilter, err := json.Marshal([]string{input.Tag})
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tag filter"})
            return
        }
        db = db.Where("tags @> ?", tagFilter)
    } else if input.Query != "" {
        db = db.Where("title ILIKE ?", "%"+input.Query+"%")
    }

    if err := db.
        Preload("Author", func(db *gorm.DB) *gorm.DB {
            return db.Select("profile_img, username, fullname")
        }).
        Order("published_at desc").
        Select("blog_id, title, des, banner, comments, total_likes, total_comments, total_reads, total_parent_comments, tags, published_at").
        Offset((input.Page - 1) * searchBlogsLimit).
        Limit(searchBlogsLimit).
        Find(&blogs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't fetch blogs"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"blogs": blogs})
}

func CreateBlog(c *gin.Context) {
    var input struct {
        Title   string   `json:"title" binding:"required"`
        Des     string   `json:"des"`
        Banner  string   `json:"banner"`
        Tags    []string `json:"tags"`
        Content struct {
            Blocks []interface{} `json:"blocks"`
        } `json:"content"`
        Draft bool `json:"draft"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    authorID, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    
    // Convert authorID to uuid.UUID
    authorUUID, err := uuid.Parse(authorID.(string))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    // Validate fields for non-draft blogs
    if !input.Draft {
        if len(input.Des) == 0 || len(input.Des) > 200 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "A blog description under 200 characters must be provided"})
            return
        }

        if len(input.Banner) == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "A blog banner must be provided"})
            return
        }

        if len(input.Content.Blocks) == 0 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Blog content must be provided"})
            return
        }

        if len(input.Tags) == 0 || len(input.Tags) > 10 {
            c.JSON(http.StatusBadRequest, gin.H{"error": "1-10 Blog tags must be provided"})
            return
        }
    }

    if len(input.Title) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "A title must be provided"})
        return
    }

    // Sanitize tags
    for i := range input.Tags {
        input.Tags[i] = strings.ToLower(input.Tags[i])
    }

    blogID := input.Title
    blogID = strings.Map(func(r rune) rune {
        if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
            return r
        }
        return ' '
    }, blogID)
    blogID = strings.Join(strings.Fields(blogID), "-") + "-" + gonanoid.Must()

    // Convert content blocks to JSON bytes
    contentBytes, err := json.Marshal(input.Content.Blocks)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process content"})
        return
    }

    // Convert tags to JSON bytes
    tagsBytes, err := json.Marshal(input.Tags)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tags"})
        return
    }
    activity := models.Activity{
        TotalLikes:           0,
        TotalComments:        0,
        TotalReads:           0,
        TotalParentComments:  0,
    }

    blog := models.Blog{
        Title:       input.Title,
        Des: input.Des,
        Banner:      input.Banner,
        Content:     contentBytes,
        Tags:        tagsBytes,
        AuthorID:    authorUUID,
        BlogID:      blogID,
        Draft:       input.Draft,
        Activity:    activity,
    }
    fmt.Println(blog)
    if err := utils.DB.Create(&blog).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create blog"})
        return
    }

    increment := 1
    if input.Draft {
        increment = 0
    }

    if err := utils.DB.Model(&models.User{}).Where("id = ?", authorUUID).Updates(map[string]interface{}{
        "total_posts": gorm.Expr("total_posts + ?", increment),
        "blogs": gorm.Expr("array_append(blogs, ?)", blog.BlogID),
    }).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update posts count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": blog.BlogID})
}
