package main

import (
	"github.com/doohee323/tz-s3-repo-mgmt/controllers"
	"github.com/doohee323/tz-s3-repo-mgmt/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

func contains(arr []string, str string) bool {
	for _, path := range arr {
		if strings.HasSuffix(str, path) {
			return true
		}
	}
	return false
}
func RedirectOnNotFound(c *gin.Context) {
	uri := c.Request.RequestURI
	allowedPaths := []string{"/", "#", ".js", ".css", ".map", ".apk", ".zip", ".woff2", ".ttf", ".png", ".gif"}
	if c.Writer.Status() == http.StatusNotFound && !contains(allowedPaths, uri) {
		log.Println("404 error")
		c.Redirect(http.StatusMovedPermanently, "/")
	}
}

func main() {
	r := gin.Default()

	models.ConnectDataBase()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		//AllowMethods:     []string{"PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	sdir := "assets"
	r.NoRoute(gin.WrapH(http.FileServer(gin.Dir(sdir, false))))
	r.Use(RedirectOnNotFound)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/logout", controllers.Logout)
		}
		user := api.Group("/user")
		{
			//utils.JwtAuthMiddleware(),
			user.PUT("", controllers.SaveUser)
			user.POST("", controllers.SaveUser)
			api.GET("/user", controllers.GetUser)
			api.DELETE("/user", controllers.DeleteUser)
			api.GET("/users", controllers.GetUsers)
		}
		s3repo := api.Group("/s3repo")
		{
			s3repo.DELETE("", controllers.DeleteS3Repo)
			api.GET("/s3repos", controllers.S3Repos)
			s3repo.GET("/download", controllers.Download)
			s3repo.POST("/upload", controllers.UploadHandler)
			s3repo.POST("/refreshCache", controllers.RefreshCache)
		}
	}

	log.Println("server started with 8080 port")
	err := r.Run(":8080")
	if err != nil {
		log.Println("failed to start with 8080 port")
		return
	}
}
