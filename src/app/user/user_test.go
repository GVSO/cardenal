package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"

	"github.com/gvso/cardenal/src/app/database/entity"
)

var linkedinToken = &oauth2.Token{
	AccessToken: "token123",
	Expiry:      time.Now().Add(60 * 24 * time.Hour),
}

var user = entity.User{
	LinkedInID:    "linkedin_id123",
	FirstName:     "John",
	LastName:      "Smith",
	LinkedInToken: *linkedinToken,
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

	userMap, err := ProcessUserAuth(data, linkedinToken)

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

	userMap, err = ProcessUserAuth(data, linkedinToken)

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

	// Saves current function and restores it at the end.
	old := insertUser
	defer func() { insertUser = old }()

	insertUserCalled := false

	insertUser = func(userArg *entity.User) (interface{}, error) {

		// Asserts that function is called with correct argument.
		assert.Equal(&user, userArg)

		result := mongo.InsertOneResult{
			InsertedID: "document123",
		}

		insertUserCalled = true

		return result, nil
	}

	data, _ := json.Marshal(user)

	userMap, err := ProcessUserAuth(data, linkedinToken)

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

	userMap, err := ProcessUserAuth(data, linkedinToken)

	expected := map[string]string{
		"linkedin_id": "linkedin_id123",
		"first_name":  "John",
		"last_name":   "Smith",
	}

	assert.True(getUserByLinkedInIDCalled)
	assert.Nil(err)
	assert.Equal(expected, userMap)
}
