package routes

import (
    "backend-go/controllers"
    "backend-go/middleware"
    "github.com/gin-gonic/gin"
)

func BlogRoutes(router *gin.Engine) {
    blogGroup := router.Group("/")
    {
        blogGroup.POST("/latest-blogs", controllers.GetLatestBlogs)
        blogGroup.POST("/latest-blogs-count", controllers.GetLatestBlogsCount)
        blogGroup.GET("/trending-blogs", controllers.GetTrendingBlogs)
        blogGroup.POST("/search-blogs-count", controllers.GetSearchBlogsCount)
        blogGroup.POST("/search-blogs", controllers.SearchBlogs)
        blogGroup.POST("/create-blog", middleware.Authenticate(), controllers.CreateBlog)

    }
}
