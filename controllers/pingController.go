package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func PingController(c *gin.Context)  {
	log.Fatal("Pong")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}