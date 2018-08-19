package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/gvso/cardenal/src/app/settings"
	"github.com/gvso/cardenal/src/app/utils/timeutils"
)

// CreateToken returns a JWT token.
var CreateToken = func(user map[string]string) (string, error) {

	ttl := time.Duration(timeutils.GetSeconds(7)) * time.Second

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"first_name":  user["first_name"],
		"last_name":   user["last_name"],
		"linkedin_id": user["linkedin_id"],
		"exp":         time.Now().UTC().Add(ttl).Unix(),
	})

	return token.SignedString(settings.JwtSecret)
}
