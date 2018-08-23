package user

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/stretchr/testify/assert"

	entity "github.com/gvso/cardenal/src/app/db/entity/user"
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
	AccessToken:   "access_token123",
}

func TestProcessUserAuth(t *testing.T) {

	assert := assert.New(t)

	data, _ := json.Marshal(user)

	testProcessUserAuthNotExistingUser(assert)

	testProcessUserAuthExistingUser(assert)

	/*****************************************************************************
	 ***************** Tests when insertUser returns an error. *******************
	 ****************************************************************************/
	old := updateUserByLinkedInID
	defer func() { updateUserByLinkedInID = old }()

	updateUserByLinkedInID = func(id string, update interface{}, fields ...string) (*entity.User, error) {
		return nil, errors.New("document does not exist")
	}

	insertUser = func(user *entity.User) (interface{}, error) {
		return nil, errors.New("document could not be inserted")
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
		return errors.New("could not unmarsh user")
	}

	userMap, err = ProcessUserAuth(data, linkedinToken)

	assert.Nil(userMap)
	assert.Equal("could not unmarsh user", err.Error())
}

func TestGetUserMap(t *testing.T) {

	assert := assert.New(t)

	// Saves current function and restores it at the end.
	old := createToken
	defer func() { createToken = old }()

	createToken = func(user map[string]string) (string, error) {
		return "access_token123", nil
	}

	userMap, err := getUserMap(&user)

	expected := map[string]string{
		"linkedin_id": "linkedin_id123",
		"first_name":  "John",
		"last_name":   "Smith",
		"token":       "access_token123",
	}

	assert.Nil(err)
	assert.Equal(expected, userMap)

	/*****************************************************************************
	 ***************** Tests when createToken returns an error. ******************
	 ****************************************************************************/

	// Overwrites createToken function.
	createToken = func(user map[string]string) (string, error) {
		return "", errors.New("could not create token")
	}

	userMap, err = getUserMap(&user)

	assert.Nil(userMap)
	assert.Equal("could not create token", err.Error())

}

// Tests when ProcessUserAuth receives a user that does not exists in database.
//
// It mocks updateUserByLinkedInID, making it return a nil user and checks that
// insertUser is called in this case.
func testProcessUserAuthNotExistingUser(assert *assert.Assertions) {

	// Saves current function and restores it at the end.
	oldGetUser := updateUserByLinkedInID
	defer func() { updateUserByLinkedInID = oldGetUser }()

	updateUserByLinkedInID = func(id string, update interface{}, fields ...string) (*entity.User, error) {
		return nil, errors.New("document does not exist")
	}

	// Saves current function and restores it at the end.
	oldCreateToken := createToken
	defer func() { createToken = oldCreateToken }()

	createToken = createTokenMock

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
		"token":       "access_token123",
	}

	assert.True(insertUserCalled)
	assert.Nil(err)
	assert.Equal(expected, userMap)
}

// Tests when ProcessUserAuth receives a user that exists in database.
//
// It mocks updateUserByLinkedInID, making it return a non-nil user and
// checks that the expected values are returned.
func testProcessUserAuthExistingUser(assert *assert.Assertions) {

	// Saves current function and restores it at the end.
	old := updateUserByLinkedInID
	defer func() { updateUserByLinkedInID = old }()

	// Saves current function and restores it at the end.
	oldCreateToken := createToken
	defer func() { createToken = oldCreateToken }()

	createToken = createTokenMock

	updateUserByLinkedInID = func(id string, update interface{}, fields ...string) (*entity.User, error) {

		// Asserts that the value of update is set correctly.
		expected := map[string]interface{}{
			"$set": map[string]string{
				"access_token": "access_token123",
			},
		}
		assert.Equal(expected, update)

		data, _ := json.Marshal(user)

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
		"token":       "access_token123",
	}

	assert.Nil(err)
	assert.Equal(expected, userMap)
}

// Mocks createToken function.
var createTokenMock = func(user map[string]string) (string, error) {
	return "access_token123", nil
}
