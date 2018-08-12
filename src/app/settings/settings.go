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

// MongoDBConfig is the structure of MongoDB settings.
type MongoDBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

// Development determines if it is development environment or not.
var Development bool

// Port is the port in which Go app is running.
var Port string

// JwtSecret is the secret for JWT token encoding.
var JwtSecret []byte

// LinkedIn holds settings for LinkedIn client.
var LinkedIn LinkedInConfig

// MongoDB holds settings for MongoDB connection
var MongoDB MongoDBConfig

// MongoPassw

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
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))

	LinkedIn.ClientID = os.Getenv("LINKEDIN_CLIENT_ID")
	LinkedIn.ClientSecret = os.Getenv("LINKEDIN_CLIENT_SECRET")
	LinkedIn.RedirectURLHost = os.Getenv("LINKEDIN_REDIRECT_URL_HOST")

	MongoDB.Host = os.Getenv("MONGO_HOST")
	MongoDB.Port = os.Getenv("MONGO_PORT")
	MongoDB.User = os.Getenv("MONGO_USER")
	MongoDB.Password = os.Getenv("MONGO_PASSWORD")
}
