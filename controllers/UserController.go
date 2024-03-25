package controllers

import (
	"github.com/doohee323/tz-s3-repo-mgmt/models"
	"github.com/doohee323/tz-s3-repo-mgmt/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type RegisterInput struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Description string `json:"description" binding:""`
	Role        string `json:"role" binding:""`
}

func GetUsers(c *gin.Context) {
	username := c.Query("username")
	log.Println(username)
	var users []models.User
	if username == "" {
		users2, err := models.GetUsers()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		users = users2
	} else {
		user, err := models.GetUserByUsername(username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	//username := c.MustGet("username").(string)
	username := c.Query("username")
	user, err := models.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username":    user.Username,
		"password":    user.Password,
		"name":        user.Name,
		"email":       user.Email,
		"description": user.Description,
		"role":        user.Role,
	})
}

func SaveUser(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		Username:    input.Username,
		Password:    input.Password,
		Name:        input.Name,
		Email:       input.Email,
		Description: input.Description,
		Role:        input.Role,
	}
	userOri, _ := models.GetUserByUsername(input.Username)
	if userOri.Username == "" {
		_, err2 := user.SaveUser()
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
	} else {
		user.ID = userOri.ID
		_, err2 := user.UpdateUser()
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err2.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	username := c.Query("username")
	userOri, err2 := models.GetUserByUsername(username)
	if err2 != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	_, err := userOri.DeleteUser()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := models.LoginCheck(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	utils.TokenCookie(c)
	user, err := models.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	err2 := godotenv.Load(".env")
	if err2 != nil {
		log.Fatalf("Error loading .env file")
	}
	s3Repos := os.Getenv("S3REPOS")
	awsRegions := os.Getenv("AWS_REGIONS")

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"config": gin.H{
			"s3Repos":    s3Repos,
			"awsRegions": awsRegions,
		},
		"user": gin.H{
			"username":    user.Username,
			"password":    user.Password,
			"name":        user.Name,
			"email":       user.Email,
			"description": user.Description,
			"role":        user.Role,
		}})
}

func Logout(c *gin.Context) {
	token := c.PostForm("token")
	if token == "" {
		c.SetCookie("session_id", "", 0, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
		return
	}
	cookies := c.Request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			if token == cookie.Value {
				token = ""
			}
			break
		}
	}
	if token == "" {
		c.SetCookie("session_id", "", 0, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully."})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Logout failed."})
	}
}
