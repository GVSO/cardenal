package main

import (
	"log"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gvso/cardenal/src/app/jwt"
	"github.com/gvso/cardenal/src/app/settings"
)

func main() {
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

	data := router.Group("/api/data").Use(jwt.Validate)
	{
		data.GET("/user", userHandler)
	}

	// Used to load css, js, and image files.
	router.Use(static.Serve("/", static.LocalFile("./client/dist", true)))

	// If route was not defined in Go server, make React handle route.
	router.NoRoute(func(c *gin.Context) {
		c.File("./client/dist/index.html")
	})

	return router
}
