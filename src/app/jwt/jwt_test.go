package jwt

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gvso/cardenal/src/app/global/mocks"
	"github.com/gvso/cardenal/src/app/settings"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {

	assert := assert.New(t)

	user := map[string]string{
		"id":         "JohnId123",
		"first_name": "John",
		"last_name":  "Smith",
	}

	token, err := CreateToken(user)

	assert.Nil(err)
	assert.True(isTokenValid(token))
}

func TestValidateHelper(t *testing.T) {

	assert := assert.New(t)

	context := &globalmocks.GinContext{}

	testValidateHelperWithNoToken(context, assert)
	testValidateHelperWithValidToken(context, assert)
	testValidateHelperWithInvalidToken(context, assert)

}

func TestIsTokenInDatabase(t *testing.T) {
	assert := assert.New(t)

	old := isUserTokenValid
	defer func() { isUserTokenValid = old }()

	isUserTokenValid = func(linkedinID string, token string) bool {
		assert.Equal("linkedin_id123", linkedinID)
		assert.Equal("token123", token)

		return true
	}

	isTokenInDatabase("linkedin_id123", "token123")
}

// Test case for when no token was provided in cookie.
func testValidateHelperWithNoToken(c *globalmocks.GinContext, assert *assert.Assertions) {
	validateHelper(c)

	assert.True(c.StringCall.Called)
	assert.Equal(403, c.StringCall.Code)
	assert.Equal("You are not authenticated", c.StringCall.Format)

	assert.True(c.AbortCall.Called)

	assert.False(c.NextCall.Called)

	// Resets values for next tests.
	resetCallValues(c)
}

// Test case for when a valid token was provided in cookie.
func testValidateHelperWithValidToken(c *globalmocks.GinContext, assert *assert.Assertions) {
	tokenString := generateToken(true)
	c.SetCookie("token", tokenString, 10, "", "", false, true)

	token, _ := jwt.ParseWithClaims(tokenString, &TokenClaims{}, KeyFunction)
	claims := token.Claims.(*TokenClaims)

	old := isTokenInDatabase
	defer func() { isTokenInDatabase = old }()

	// Tests when token is the one in database.
	isTokenInDatabase = func(linkedin_id string, token string) bool {
		assert.Equal(claims.LinkedInID, linkedin_id)
		assert.Equal(tokenString, token)

		return true
	}

	validateHelper(c)

	assert.True(c.SetCall.Called)
	assert.Equal("token", c.SetCall.Key)
	assert.Equal(token.Claims, c.SetCall.Value)

	assert.True(c.NextCall.Called)

	assert.False(c.AbortCall.Called)

	// Resets values for next tests.
	resetCallValues(c)

	// Tests when token is not the one in database.
	isTokenInDatabase = func(linkedin_id string, token string) bool {
		assert.Equal(claims.LinkedInID, linkedin_id)
		assert.Equal(tokenString, token)

		return false
	}

	validateHelper(c)

	assert.False(c.SetCall.Called)
	assert.False(c.NextCall.Called)

	assert.True(c.AbortCall.Called)
	assert.True(c.StringCall.Called)
	assert.Equal(403, c.StringCall.Code)
	assert.Equal("You are not authenticated", c.StringCall.Format)

	// Resets values for next tests.
	resetCallValues(c)
}

// Test case for when a invalid token was provided in cookie.
func testValidateHelperWithInvalidToken(c *globalmocks.GinContext, assert *assert.Assertions) {
	settings.Development = true

	tokenString := generateToken(false)
	c.SetCookie("token", tokenString, 10, "", "", false, true)

	validateHelper(c)

	assert.True(c.StringCall.Called)
	assert.Equal(500, c.StringCall.Code)
	assert.Equal("Invalid token", c.StringCall.Format)

	assert.False(c.SetCall.Called)
	assert.True(c.AbortCall.Called)
	assert.False(c.NextCall.Called)

	// Resets values for next tests.
	resetCallValues(c)
}

// Checks that generated token stored in cookie is valid.
//
// It parses the token string and check for errors and validity of token. If no
// errors are produced and token is valid, return true.
func isTokenValid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, KeyFunction)

	if err != nil {
		return false
	}

	if token.Valid {
		return true
	}

	return false
}

// Generates a token
//
// It generates a valid or invalid token based on the value passed as an
// argument
func generateToken(valid bool) string {

	if !valid {
		return "AnInvalidToken"
	}

	ttl := 10 * time.Second

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"first_name":  "John",
		"last_name":   "Smith",
		"linkedin_id": "linkedin_id123",
		"exp":         time.Now().UTC().Add(ttl).Unix(),
	})

	tokenString, _ := token.SignedString(settings.JwtSecret)

	return tokenString
}

func resetCallValues(c *globalmocks.GinContext) {
	// Resets values for next tests.
	c.SetCall.Called = false
	c.NextCall.Called = false
	c.AbortCall.Called = false
	c.StringCall.Called = false
}
