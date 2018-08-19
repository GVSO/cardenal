package user

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gvso/cardenal/src/app/database/entity"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/stretchr/testify/assert"
)

var user = entity.User{
	LinkedInID: "linkedin_id123",
	FirstName:  "John",
	LastName:   "Smith",
}

func TestProcessUserAuth(t *testing.T) {

	assert := assert.New(t)

	data, _ := json.Marshal(user)

	// Saves current function and restores it at the end.
	old := insertUser
	defer func() { insertUser = old }()

	insertUser = func(user *entity.User) (interface{}, error) {
		result := mongo.InsertOneResult{
			InsertedID: "document123",
		}

		return result, nil
	}

	userMap, err := ProcessUserAuth(data)

	expected := map[string]string{
		"linkedin_id": "linkedin_id123",
		"first_name":  "John",
		"last_name":   "Smith",
	}

	assert.Nil(err)
	assert.Equal(expected, userMap)

	// Tests when insertUser returns an error.
	insertUser = func(user *entity.User) (interface{}, error) {
		return nil, fmt.Errorf("document could not be inserted")
	}

	userMap, err = ProcessUserAuth(data)

	assert.Nil(userMap)
	assert.Equal("document could not be inserted", err.Error())

	// Tests when jsonUnmarshal returns an error.

	// Saves current function and restores it at the end.
	oldJSONUnmarshal := jsonUnmarshal
	defer func() { jsonUnmarshal = oldJSONUnmarshal }()

	jsonUnmarshal = func(data []byte, v interface{}) error {
		return fmt.Errorf("could not unmarsh user")
	}

	userMap, err = ProcessUserAuth(data)

	assert.Nil(userMap)
	assert.Equal("could not unmarsh user", err.Error())
}
