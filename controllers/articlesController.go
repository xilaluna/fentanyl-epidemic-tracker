package controllers

import (
	"github.com/gin-gonic/gin"
)

func GetArticles(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "articles",
	})
}