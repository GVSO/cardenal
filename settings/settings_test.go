package settings

import (
	"fmt"
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
				"LINKEDIN_CLIENT_ID":         "client123",
				"LINKEDIN_CLIENT_SECRET":     "secret123",
				"LINKEDIN_REDIRECT_URL_HOST": "http://localhost",
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
				assert.Equal("", Port, "port in settings to be empty")

				return
			}

			// Assert Development value.
			expected, _ := strconv.ParseBool(environment["DEVELOPMENT"])
			msg := fmt.Sprintf("settings.Development should be %v; got %v", environment["DEVELOPMENT"], Development)
			assert.Equal(expected, Development, msg)

			// Assert Port value.
			msg = fmt.Sprintf("settings.Port should be %v; got %v", environment["GO_PORT"], Port)
			assert.Equal(environment["GO_PORT"], Port, msg)

			// Assert LinkedIn.ClientID value.
			msg = fmt.Sprintf("settings.LinkedIn.ClientID should be %v; got %v", environment["LINKEDIN_CLIENT_ID"], LinkedIn.ClientID)
			assert.Equal(environment["LINKEDIN_CLIENT_ID"], LinkedIn.ClientID, msg)

			// Assert LinkedIn.ClientSecret value.
			msg = fmt.Sprintf("settings.LinkedIn.ClientSecret should be %v; got %v", environment["LINKEDIN_CLIENT_SECRET"], LinkedIn.ClientSecret)
			assert.Equal(environment["LINKEDIN_CLIENT_SECRET"], LinkedIn.ClientSecret, msg)

			// Assert LinkedIn.RedirectURLHost value.
			msg = fmt.Sprintf("settings.LinkedIn.RedirectURLHost should be %v; got %v", environment["LINKEDIN_REDIRECT_URL_HOST"], LinkedIn.RedirectURLHost)
			assert.Equal(environment["LINKEDIN_REDIRECT_URL_HOST"], LinkedIn.RedirectURLHost, msg)

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
