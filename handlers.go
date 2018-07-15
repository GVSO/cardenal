package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gvso/cardenal/linkedin"
)

func loginHandler(c *gin.Context) {
	linkedin.Login(c)
}

func callbackHandler(c *gin.Context) {
	linkedin.Callback(c)
}
