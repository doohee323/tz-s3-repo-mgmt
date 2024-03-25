package controllers

import (
	"fmt"
	"github.com/doohee323/tz-s3-repo-mgmt/models"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"os"
)

func S3Repos(c *gin.Context) {
	awsRegion := c.Query("awsRegion")
	s3repo := c.Query("s3repo")
	outStr, err := models.GetS3Repos(awsRegion, s3repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": outStr})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": outStr})
	}
}

func DeleteS3Repo(c *gin.Context) {
	s3repo := c.Query("s3repo")
	object := c.Query("object")
	_, err := models.DeleteS3Repo(s3repo, object)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "S3Repo deleted successfully"})
}

func Download(c *gin.Context) {
	s3repo := c.Query("s3repo")
	object := c.Query("object")
	_, err := models.GetS3Repo(s3repo, object)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, err := os.Open(object)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer file.Close()
	c.Header("Content-Disposition", "attachment; filename="+object)
	c.Header("Content-Type", "application/octet-stream")
	c.File(object)
}

type Upload struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func UploadHandler(c *gin.Context) {
	s3repo := c.PostForm("s3repo")
	if s3repo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "s3repo is required."})
		return
	}
	var upload Upload
	if err := c.ShouldBind(&upload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, _ := upload.File.Open()
	sourceFileName, err := models.WriteUploadFile(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err2 := models.UploadS3Repo(s3repo, sourceFileName, upload.File.Filename)
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Uploaded file %s successfully.", "")})
}
