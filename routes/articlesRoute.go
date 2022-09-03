package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/xilaluna/fentanyl-epidemic-tracker/controllers"
)

func ArticlesRoute(router *gin.Engine) {
	router.GET("/articles", controllers.GetArticles)
}