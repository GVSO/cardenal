package user

import (
	"encoding/json"

	"github.com/gvso/cardenal/src/app/database"
)

// ProcessUserAuth handles user authentication/registration after user
// has authenticated on LinkedIn.
func ProcessUserAuth(user []byte) (map[string]string, error) {
	userMap := make(map[string]interface{})

	err := json.Unmarshal(user, &userMap)
	if err != nil {
		panic(err)
	}

	database.InsertUser(userMap)

	return getUserData(userMap), nil
}

func getUserData(userMap map[string]interface{}) map[string]string {

	return map[string]string{
		"id":         userMap["id"].(string),
		"first_name": userMap["firstName"].(string),
		"last_name":  userMap["lastName"].(string),
	}
}
