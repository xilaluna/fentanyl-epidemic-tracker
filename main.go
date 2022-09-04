package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xilaluna/fentanyl-epidemic-tracker/configs"
	"github.com/xilaluna/fentanyl-epidemic-tracker/routes"
)


func main() {
	router := gin.Default()

	configs.ConnectDB()

	router.StaticFile("/", "./static/index.html")
	routes.ArticlesRoute(router)
	routes.PingRoute(router)
	routes.ScrapeRoute(router)

	defer configs.CloseDB()
	router.Run("0.0.0.0:" + os.Getenv("PORT"))
}