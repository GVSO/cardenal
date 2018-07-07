package settings

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LinkedInConfig is the structure of LinkedIn settings.
type LinkedInConfig struct {
	ClientID        string
	ClientSecret    string
	RedirectURLHost string
}

// Port is the port in which Go app is running.
var Port string

// Development determines if it is development environment or not.
var Development bool

// LinkedIn holds settings for LinkedIn client.
var LinkedIn LinkedInConfig

func init() {
	godotenv.Load()

	initSettings()
}

func initSettings() {
	if os.Getenv("DEVELOPMENT") == "" {
		log.Println("Missing environment variable file.")

		return
	}

	Port = os.Getenv("GO_PORT")
	Development, _ = strconv.ParseBool(os.Getenv("DEVELOPMENT"))
	LinkedIn.ClientID = os.Getenv("LINKEDIN_CLIENT_ID")
	LinkedIn.ClientSecret = os.Getenv("LINKEDIN_CLIENT_SECRET")
	LinkedIn.RedirectURLHost = os.Getenv("LINKEDIN_REDIRECT_URL_HOST")
}
