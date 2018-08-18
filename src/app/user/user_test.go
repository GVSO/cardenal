package user

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/stretchr/testify/assert"
)

var userMap = map[string]interface{}{
	"id":         "id123",
	"firstName":  "John",
	"lastName":   "Smith",
	"extraField": "extra field",
}

func TestProcessUserAuth(t *testing.T) {

	assert := assert.New(t)

	user, _ := json.Marshal(userMap)

	// Saves current function and restores it at the end.
	old := insertUser
	defer func() { insertUser = old }()

	insertUser = func(user map[string]interface{}) (interface{}, error) {
		result := mongo.InsertOneResult{
			InsertedID: "document123",
		}

		return result, nil
	}

	userMap, err := ProcessUserAuth(user)

	expected := map[string]string{
		"id":         "id123",
		"first_name": "John",
		"last_name":  "Smith",
	}

	assert.Nil(err)
	assert.Equal(expected, userMap)

	// Tests when insertUser returns an error.
	insertUser = func(user map[string]interface{}) (interface{}, error) {
		return nil, fmt.Errorf("document could not be inserted")
	}

	userMap, err = ProcessUserAuth(user)

	assert.Nil(userMap)
	assert.Equal("document could not be inserted", err.Error())

	// Tests when jsonUnmarshal returns an error.

	// Saves current function and restores it at the end.
	oldJSONUnmarshal := jsonUnmarshal
	defer func() { jsonUnmarshal = oldJSONUnmarshal }()

	jsonUnmarshal = func(data []byte, v interface{}) error {
		return fmt.Errorf("could not unmarsh user")
	}

	userMap, err = ProcessUserAuth(user)

	assert.Nil(userMap)
	assert.Equal("could not unmarsh user", err.Error())
}

func TestParseUserData(t *testing.T) {

	assert := assert.New(t)

	user := parseUserData(userMap)

	expected := map[string]interface{}{
		"id":         "id123",
		"first_name": "John",
		"last_name":  "Smith",
	}

	assert.Equal(expected, user)
}
