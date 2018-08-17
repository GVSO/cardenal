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
	Database string
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

func init() {
	godotenv.Load()

	initSettings()
}

var initSettings = func() {
	if os.Getenv("DEVELOPMENT") == "" {
		log.Println("Missing environment variable file.")

		return
	}

	Development, _ = strconv.ParseBool(os.Getenv("DEVELOPMENT"))
	Port = os.Getenv("GO_PORT")
	JwtSecret = []byte(os.Getenv("JWT_SECRET"))

	LinkedIn = LinkedInConfig{
		os.Getenv("LINKEDIN_CLIENT_ID"),
		os.Getenv("LINKEDIN_CLIENT_SECRET"),
		os.Getenv("LINKEDIN_REDIRECT_URL_HOST"),
	}

	MongoDB = MongoDBConfig{
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_PORT"),
		os.Getenv("MONGO_DB"),
		os.Getenv("MONGO_USER"),
		os.Getenv("MONGO_PASSWORD"),
	}
}
