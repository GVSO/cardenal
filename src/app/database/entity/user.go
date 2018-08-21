package entity

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"golang.org/x/oauth2"
)

// User is the struct for a user.
type User struct {
	ID            objectid.ObjectID `json:"-" bson:"_id"`
	LinkedInID    string            `json:"id" bson:"linkedin_id"`
	FirstName     string            `json:"firstName" bson:"first_name"`
	LastName      string            `json:"lastName" bson:"last_name"`
	LinkedInToken oauth2.Token      `json:"-" bson:"linkedin_token"`
}
