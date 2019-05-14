package user

import (
	"encoding/json"

	"golang.org/x/oauth2"

	entity "github.com/gvso/cardenal/src/app/db/entity/user"
	"github.com/gvso/cardenal/src/app/jwt"
)

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as pointers to
 * function. The advantage of doing this is that these variables can later be
 * overwritten in the testing files.
 */
var createTokenString = jwt.CreateTokenString
var updateUserByLinkedInID = entity.UpdateUserByLinkedInID
var insertUser = entity.InsertUser
var jsonUnmarshal = json.Unmarshal

// ProcessUserAuth handles user authentication/registration after user
// has authenticated on LinkedIn.
var ProcessUserAuth = func(data []byte,
	linkedinToken *oauth2.Token) (map[string]string, error) {

	user := &entity.User{}

	err := jsonUnmarshal(data, user)
	if err != nil {
		return nil, err
	}

	userMap, err := getUserMap(user)

	// Tries to update user access token in database.
	update := map[string]interface{}{
		"$set": map[string]string{
			"access_token": userMap["token"],
		},
	}

	updatedUser, _ := updateUserByLinkedInID(user.LinkedInID, update)

	// If user was not found, creates new user.
	if updatedUser == nil {

		user.LinkedInToken = *linkedinToken
		user.AccessToken = userMap["token"]

		_, err = insertUser(user)
		if err != nil {
			return nil, err
		}

	}

	return userMap, nil

}

var getUserMap = func(user *entity.User) (map[string]string, error) {

	userMap := map[string]string{
		"linkedin_id": user.LinkedInID,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
	}

	token, err := createTokenString(userMap)
	if err != nil {
		return nil, err
	}

	userMap["token"] = token

	return userMap, nil
}
