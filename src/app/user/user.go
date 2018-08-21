package user

import (
	"encoding/json"

	"golang.org/x/oauth2"

	"github.com/gvso/cardenal/src/app/database"
	"github.com/gvso/cardenal/src/app/database/entity"
)

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as pointers to
 * function. The advantage of doing this is that these variables can later be
 * overwritten in the testing files.
 */
var getUserByLinkedInID = database.GetUserByLinkedInID
var insertUser = database.InsertUser
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

	exists, _ := getUserByLinkedInID(user.LinkedInID)

	// If user was not found, creates new user.
	if exists == nil {

		user.LinkedInToken = *linkedinToken

		_, err = insertUser(user)
		if err != nil {
			return nil, err
		}

	}

	return map[string]string{
		"linkedin_id": user.LinkedInID,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
	}, nil

}
