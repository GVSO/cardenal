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

// Development determines if it is development environment or not.
var Development bool

// Port is the port in which Go app is running.
var Port string

// JwtKey is the secret key for JWT token encryption.
var JwtKey []byte

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

	Development, _ = strconv.ParseBool(os.Getenv("DEVELOPMENT"))
	Port = os.Getenv("GO_PORT")
	JwtKey = []byte(os.Getenv("JWT_KEY"))

	LinkedIn.ClientID = os.Getenv("LINKEDIN_CLIENT_ID")
	LinkedIn.ClientSecret = os.Getenv("LINKEDIN_CLIENT_SECRET")
	LinkedIn.RedirectURLHost = os.Getenv("LINKEDIN_REDIRECT_URL_HOST")
}