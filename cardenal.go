package main

import (
	"log"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gvso/cardenal/settings"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := setupRouter()

	port := ":" + settings.Port

	log.Print("Server application started at ", "http://localhost"+port)

	router.Run(port)
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	services := router.Group("/api/services")
	{
		services.GET("/login", loginHandler)
		services.GET("/login/callback", callbackHandler)
	}

	router.Use(static.Serve("/", static.LocalFile("./client/dist", true)))

	return router
}
