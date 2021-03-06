package user

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"

	"github.com/gvso/cardenal/src/app/db"
)

// User is the struct for a user.
type User struct {
	ID            primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	LinkedInID    string             `json:"id" bson:"linkedin_id"`
	FirstName     string             `json:"localizedFirstName" bson:"first_name"`
	LastName      string             `json:"localizedLastName" bson:"last_name"`
	AccessToken   string             `json:"-" bson:"access_token"`
	LinkedInToken oauth2.Token       `json:"-" bson:"linkedin_token"`
}

var collection db.MongoCollection

func init() {
	collection = db.GetCollection("users")
}

// InsertUser inserts a new user.
var InsertUser = func(user *User) (interface{}, error) {

	res, err := collection.InsertOne(nil, user)
	if err != nil {
		return nil, err
	}

	id := res.InsertedID

	return id, nil
}

// GetUserByLinkedInID returns the document containing the user with the
// provided ID.
//
// If no user with such ID exists in database, an error and a nil user are
// returned.
var GetUserByLinkedInID = func(id string, fields ...string) (*User, error) {

	user := User{}
	filter := bson.D{{"linkedin_id", id}}

	err := findOne(filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserByLinkedInID updates and returns the document containing the user
// with the provided ID.
//
// If no user with such ID exists in database, an error and a nil user are
// returned.
var UpdateUserByLinkedInID = func(id string, update interface{},
	fields ...string) (*User, error) {

	user := User{}
	filter := bson.D{{"linkedin_id", id}}

	err := findOneAndUpdate(&filter, update).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// IsTokenValid determines if token is valid for the given user.
//
// It checks that the token passed to the function is the current one in the
// database.
var IsTokenValid = func(linkedinID string, token string) bool {
	user := User{}
	filter := map[string]string{
		"linkedin_id":  linkedinID,
		"access_token": token,
	}

	err := findOne(filter).Decode(&user)
	if err != nil {
		return false
	}

	return true
}

// Returns a SingleResult that meets the filter criteria.
var findOne = func(filter interface{}) db.SingleResult {

	return collection.FindOne(nil, filter)
}

// Updates and returns a SingleResult that meets the filter criteria.
var findOneAndUpdate = func(filter *bson.D,
	update interface{}) db.SingleResult {

	return collection.FindOneAndUpdate(nil, filter, update)
}
