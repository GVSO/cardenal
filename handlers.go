package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gvso/cardenal/src/app/database"

	"github.com/gvso/cardenal/src/app/linkedin"
)

func loginHandler(c *gin.Context) {
	linkedin.Login(c)
}

func callbackHandler(c *gin.Context) {
	linkedin.Callback(c)
}

func userHandler(c *gin.Context) {
	c.String(200, "success")

	database.GetUserByLinkedInID("MrZtdNnOwu")
}
