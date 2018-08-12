package settings

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitiSettings(t *testing.T) {

	tests := map[string]struct {
		environment map[string]string
	}{
		"environment variables not provided.": {
			environment: map[string]string{},
		},
		"environment variables exist.": {
			environment: map[string]string{
				"DEVELOPMENT":                "true",
				"GO_PORT":                    "6000",
				"JWT_KEY":                    "key",
				"LINKEDIN_CLIENT_ID":         "client123",
				"LINKEDIN_CLIENT_SECRET":     "secret123",
				"LINKEDIN_REDIRECT_URL_HOST": "http://localhost",
				"MONGO_HOST":                 "localhost",
				"MONGO_PORT":                 "27017",
				"MONGO_DB":                   "mydb",
				"MONGO_USER":                 "mongouser",
				"MONGO_PASSWORD":             "password123",
			},
		},
	}

	for name, test := range tests {

		t.Run(name, func(t *testing.T) {
			cleanUp()

			environment := test.environment
			for variable, value := range environment {
				os.Setenv(variable, value)
			}

			initSettings()

			assert := assert.New(t)

			// No environment variables, values should be empty.
			if environment == nil {
				assert.Equal("", Port)

				return
			}

			// Asserts Development value.
			expected, _ := strconv.ParseBool(environment["DEVELOPMENT"])
			assert.Equal(expected, Development)

			// Asserts Port value.
			assert.Equal(environment["GO_PORT"], Port)

			// Asserts JWT secret value.
			assert.Equal(environment["JWT_SECRET"], string(JwtSecret))

			// Asserts LinkedIn setting values.
			assert.Equal(environment["LINKEDIN_CLIENT_ID"], LinkedIn.ClientID)
			assert.Equal(environment["LINKEDIN_CLIENT_SECRET"], LinkedIn.ClientSecret)
			assert.Equal(environment["LINKEDIN_REDIRECT_URL_HOST"], LinkedIn.RedirectURLHost)

			// Asserts MongoDB setting values.
			assert.Equal(environment["MONGO_HOST"], MongoDB.Host)
			assert.Equal(environment["MONGO_PORT"], MongoDB.Port)
			assert.Equal(environment["MONGO_DB"], MongoDB.Database)
			assert.Equal(environment["MONGO_USER"], MongoDB.User)
			assert.Equal(environment["MONGO_PASSWORD"], MongoDB.Password)
		})
	}
}

func cleanUp() {
	unsetVariables()
	unsetSettings()
}

// Unsets environment variables in case they are coming from app.
func unsetVariables() {
	variables := []string{
		"DEVELOPMENT",
		"GO_PORT",
		"LINKEDIN_CLIENT_ID",
		"LINKEDIN_CLIENT_SECRET",
		"LINKEDIN_REDIRECT_URL_HOST",
	}

	for _, variable := range variables {
		os.Unsetenv(variable)
	}
}

// Reinitializes setting variables
func unsetSettings() {
	Port = ""
	LinkedIn.ClientID = ""
	LinkedIn.ClientSecret = ""
	LinkedIn.RedirectURLHost = ""
}
