package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xilaluna/fentanyl-epidemic-tracker/controllers"
)

func ScrapeRoute(router *gin.Engine) {
	router.GET("/scrape", controllers.ScrapeController)
}