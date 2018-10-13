package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/gvso/cardenal/src/app/db/entity/user"
	"github.com/gvso/cardenal/src/app/global"
	"github.com/gvso/cardenal/src/app/settings"
)

// TokenClaims is the struct of the claims in the JWT token.
type TokenClaims struct {
	*jwt.StandardClaims
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	LinkedInID string `json:"linkedin_id"`
}

var key []byte

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as pointers to
 * function. The advantage of doing this is that these variables can later be
 * overwritten in the testing files.
 */
var isUserTokenValid = user.IsTokenValid

// Validate checks that JWT token is valid.
var Validate = func(c *gin.Context) {
	validateHelper(c)
}

// KeyFunction returns the encoding secret key.
var KeyFunction = func(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("There was an error")
	}

	return settings.JwtSecret, nil
}

// Helper function for Validate function.
var validateHelper = func(c global.GinContext) {

	tokenString, err := c.Cookie("token")

	if tokenString != "" && err == nil {
		token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, KeyFunction)

		if err != nil {
			c.Abort()

			c.String(500, "Invalid token")

			if settings.Development {
				fmt.Println(err.Error())
			}

			return
		}

		if token.Valid {

			claims := token.Claims.(*TokenClaims)
			linkedinID := (*claims).LinkedInID

			if isTokenInDatabase(linkedinID, tokenString) {
				c.Set("token", token.Claims)

				// Calls next request handler.
				c.Next()

				return
			}

		}
	}

	c.String(403, "You are not authenticated")

	c.Abort()
}

var isTokenInDatabase = func(linkedinID string, token string) bool {
	return isUserTokenValid(linkedinID, token)
}

// @TODO: Update LinkedIn Access Token.
/*var updateAccessToken = func(linkedinID string) jwt.Token {

}*/
