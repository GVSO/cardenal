package user

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"golang.org/x/oauth2"

	"github.com/gvso/cardenal/src/app/db"
)

// User is the struct for a user.
type User struct {
	ID            objectid.ObjectID `json:"-" bson:"_id"`
	LinkedInID    string            `json:"id" bson:"linkedin_id"`
	FirstName     string            `json:"firstName" bson:"first_name"`
	LastName      string            `json:"lastName" bson:"last_name"`
	LinkedInToken oauth2.Token      `json:"-" bson:"linkedin_token"`
	AccessToken   string            `json:"-" bson:"access_token"`
}

var collection db.MongoCollection

func init() {
	collection = db.GetCollection("users")
}

// InsertUser inserts a new user.
var InsertUser = func(user *User) (interface{}, error) {

	// Sets the _id field value.
	(*user).ID = objectid.New()

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
	filter := bson.NewDocument(bson.EC.String("linkedin_id", id))

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
	filter := bson.NewDocument(bson.EC.String("linkedin_id", id))

	err := findOneAndUpdate(filter, update).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Returns a single DocumentResult that meets the filter criteria.
var findOne = func(filter *bson.Document) db.DocumentResult {

	return collection.FindOne(nil, filter)
}

// Updates and returns a single DocumentResult that meets the filter criteria.
var findOneAndUpdate = func(filter *bson.Document,
	update interface{}) db.DocumentResult {

	return collection.FindOneAndUpdate(nil, filter, update)
}
