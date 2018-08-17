package user

import (
	"encoding/json"

	"github.com/gvso/cardenal/src/app/database"
)

// ProcessUserAuth handles user authentication/registration after user
// has authenticated on LinkedIn.
var ProcessUserAuth = func(user []byte) (map[string]string, error) {
	userMap := make(map[string]interface{})

	err := json.Unmarshal(user, &userMap)
	if err != nil {
		return nil, err
	}

	userMap = parseUserData(userMap)

	_, err = database.InsertUser(userMap)
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
