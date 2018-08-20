package user

import (
	"encoding/json"
	"errors"
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

	testProcessUserAuthNotExistingUser(assert)

	testProcessUserAuthExistingUser(assert)

	/*****************************************************************************
	 ***************** Tests when insertUser returns an error. *******************
	 ****************************************************************************/
	old := getUserByLinkedInID
	defer func() { getUserByLinkedInID = old }()

	getUserByLinkedInID = func(id string, fields ...string) (*entity.User, error) {
		return nil, errors.New("document does not exist")
	}

	insertUser = func(user *entity.User) (interface{}, error) {
		return nil, fmt.Errorf("document could not be inserted")
	}

	userMap, err := ProcessUserAuth(data)

	assert.Nil(userMap)
	assert.Equal("document could not be inserted", err.Error())

	/*****************************************************************************
	 *************** Tests when jsonUnmarshal returns an error. ******************
	 ****************************************************************************/

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

// Tests when ProcessUserAuth receives a user that does not exists in database.
//
// It mocks getUserByLinkedInID, making it return a nil user and checks that
// insertUser is called in this case.
func testProcessUserAuthNotExistingUser(assert *assert.Assertions) {

	// Saves current function and restores it at the end.
	oldGetUser := getUserByLinkedInID
	defer func() { getUserByLinkedInID = oldGetUser }()

	getUserByLinkedInID = func(id string, fields ...string) (*entity.User, error) {
		return nil, errors.New("document does not exist")
	}

	data, _ := json.Marshal(user)

	// Saves current function and restores it at the end.
	old := insertUser
	defer func() { insertUser = old }()

	insertUserCalled := false

	insertUser = func(user *entity.User) (interface{}, error) {

		expected := &entity.User{}
		json.Unmarshal(data, expected)

		// Asserts that function is called with correct argument.
		assert.Equal(expected, user)

		result := mongo.InsertOneResult{
			InsertedID: "document123",
		}

		insertUserCalled = true

		return result, nil
	}

	userMap, err := ProcessUserAuth(data)

	expected := map[string]string{
		"linkedin_id": "linkedin_id123",
		"first_name":  "John",
		"last_name":   "Smith",
	}

	assert.True(insertUserCalled)
	assert.Nil(err)
	assert.Equal(expected, userMap)
}

// Tests when ProcessUserAuth receives a user that does exists in database.
//
// It mocks getUserByLinkedInID, making it return a non-nil user and checks that
// getUserByLinkedInID is called in this case.
func testProcessUserAuthExistingUser(assert *assert.Assertions) {

	// Saves current function and restores it at the end.
	old := getUserByLinkedInID
	defer func() { getUserByLinkedInID = old }()

	getUserByLinkedInIDCalled := false

	getUserByLinkedInID = func(id string, fields ...string) (*entity.User, error) {

		data, _ := json.Marshal(user)

		getUserByLinkedInIDCalled = true

		user := &entity.User{}
		json.Unmarshal(data, user)

		return user, nil
	}

	data, _ := json.Marshal(user)

	userMap, err := ProcessUserAuth(data)

	expected := map[string]string{
		"linkedin_id": "linkedin_id123",
		"first_name":  "John",
		"last_name":   "Smith",
	}

	assert.True(getUserByLinkedInIDCalled)
	assert.Nil(err)
	assert.Equal(expected, userMap)
}
