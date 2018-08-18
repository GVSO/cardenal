package user

import (
	"encoding/json"

	"github.com/gvso/cardenal/src/app/database"
)

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as pointers to
 * function. The advantage of doing this is that these variables can later be
 * overwritten in the testing files.
 */
var insertUser = database.InsertUser
var jsonUnmarshal = json.Unmarshal

// ProcessUserAuth handles user authentication/registration after user
// has authenticated on LinkedIn.
var ProcessUserAuth = func(user []byte) (map[string]string, error) {
	userMap := make(map[string]interface{})

	err := jsonUnmarshal(user, &userMap)
	if err != nil {
		return nil, err
	}

	userMap = parseUserData(userMap)

	_, err = insertUser(userMap)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"id":         userMap["id"].(string),
		"first_name": userMap["first_name"].(string),
		"last_name":  userMap["last_name"].(string),
	}, nil
}

var parseUserData = func(userMap map[string]interface{}) map[string]interface{} {

	return map[string]interface{}{
		"id":         userMap["id"],
		"first_name": userMap["firstName"],
		"last_name":  userMap["lastName"],
	}
}
