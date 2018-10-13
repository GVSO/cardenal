package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gvso/cardenal/src/app/linkedin"
)

func loginHandler(c *gin.Context) {
	linkedin.Login(c)
}

func callbackHandler(c *gin.Context) {
	linkedin.Callback(c)
}

func userHandler(c *gin.Context) {
	token, _ := c.Cookie("token")
	c.String(200, token)
}
