package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xilaluna/fentanyl-epidemic-tracker/controllers"
)

func PingRoute(router *gin.Engine) {
	router.GET("/ping", controllers.PingController)
}