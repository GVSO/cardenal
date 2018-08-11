package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/gvso/cardenal/src/app/global"
	"github.com/gvso/cardenal/src/app/settings"
)

var key []byte

// Validate checks that JWT token is valid.
func Validate(c *gin.Context) {
	validateHelper(c)
}

// Helper function for Validate function.
func validateHelper(c global.GinContext) {

	authorizationValue, err := c.Cookie("token")

	if authorizationValue != "" && err == nil {
		token, err := jwt.Parse(authorizationValue, KeyFunction)

		if err != nil {
			c.Abort()

			c.String(500, "You could not be authenticated")

			if settings.Development {
				fmt.Println(err.Error())
			}

			return
		}

		if token.Valid {
			c.Set("token", token.Claims)

			// Calls next request handler.
			c.Next()

			return
		}
	}

	c.String(403, "You are not authenticated")

	c.Abort()

}

// KeyFunction returns the encoding secret key.
func KeyFunction(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("There was an error")
	}

	return settings.JwtSecret, nil
}
