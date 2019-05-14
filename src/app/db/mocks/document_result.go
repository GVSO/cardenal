package mocks

import (
	"errors"
	"reflect"
)

// SingleResult is the mock structure for database.SingleResult.
type SingleResult struct {
	DecodeCall decode
}

// Decode mocks a call to Decode.
func (_m *SingleResult) Decode(v interface{}) error {

	times := &_m.DecodeCall.times

	// Error on second and fourth call.
	if *times == 1 || *times == 3 {
		(*times)++

		_m.DecodeCall = decode{*times, true, v}

		return errors.New("could not decode document")
	}

	(*times)++

	user := reflect.ValueOf(v).Elem()

	_m.DecodeCall = decode{*times, true, user.Interface()}

	user.Field(1).SetString("linkedin_id123") // LinkedInID
	user.Field(2).SetString("John")           // FirstName
	user.Field(3).SetString("Smith")          // LastName

	return nil
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type decode struct {
	times int

	Called bool
	Value  interface{}
}
