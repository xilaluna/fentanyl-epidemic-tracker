package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xilaluna/fentanyl-epidemic-tracker/configs"
	"github.com/xilaluna/fentanyl-epidemic-tracker/controllers"
)


func main() {
	router := gin.Default()

	configs.ConnectDB()

	router.StaticFile("/", "./static/index.html")
	router.GET("/articles", controllers.GetArticles)
	router.GET("/ping", controllers.PingController)
	router.GET("/scrape", controllers.ScrapeController)

	defer configs.CloseDB()
	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}