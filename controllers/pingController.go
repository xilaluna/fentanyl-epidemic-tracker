package controllers

import "github.com/gin-gonic/gin"

func PingController(c *gin.Context)  {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}