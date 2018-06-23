package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type settings struct {
	Port        string
	Development bool
	LinkedIn    struct {
		ClientID        string
		ClientSecret    string
		RedirectURLHost string
	}
}

// Settings global settings variable.
var Settings settings

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file in root project")
	}

	Settings.Port = os.Getenv("GO_PORT")
	Settings.Development, _ = strconv.ParseBool(os.Getenv("DEVELOPMENT"))
	Settings.LinkedIn.ClientID = os.Getenv("LINKEDIN_CLIENT_ID")
	Settings.LinkedIn.ClientSecret = os.Getenv("LINKEDIN_CLIENT_SECRET")
	Settings.LinkedIn.RedirectURLHost = os.Getenv("LINKEDIN_REDIRECT_URL_HOST")
}
