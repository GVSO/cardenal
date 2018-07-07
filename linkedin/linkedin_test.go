package linkedin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gvso/cardenal/constants"
	"golang.org/x/oauth2/linkedin"

	"github.com/gvso/cardenal/settings"
	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	settings.Development = true

	assert := assert.New(t)

	client := HTTPClientMock{}

	// First time should be 200.
	data, err := getProfile(client)

	assert.Nil(err)
	assert.Equal("Successful data", string(data), "received data is not correct")

	// Second time should be 403.
	data, err = getProfile(client)

	assert.NotNil(err)
	assert.Equal("request was not successful", err.Error(), "error messages do not match")
	assert.Nil(data)

	// Third time should be 400.
	data, err = getProfile(client)

	assert.NotNil(err)
	assert.Equal("request was not successful", err.Error(), "error messages do not match")
	assert.Nil(data)

	// Fourth time should be an error when making request.
	data, err = getProfile(client)

	assert.NotNil(err)
	assert.Equal("Error on request", err.Error(), "error messages do not match")
	assert.Nil(data)
}
func TestGetConfig(t *testing.T) {

	// Override settings' values.
	settings.LinkedIn = settings.LinkedInConfig{
		ClientID:        "client123",
		ClientSecret:    "secret123",
		RedirectURLHost: "http://localhost",
	}
	settings.Port = "8000"

	config := getConfig()

	assert := assert.New(t)

	assert.Equal(settings.LinkedIn.ClientID, config.ClientID, "client ID is not correct")
	assert.Equal(settings.LinkedIn.ClientSecret, config.ClientSecret, "client secret is not correct")

	redirectURL := settings.LinkedIn.RedirectURLHost + ":" + settings.Port + constants.LinkedInRedirectPath
	assert.Equal(redirectURL, config.RedirectURL, "redirect URL is not correct")

	assert.Equal([]string{"r_basicprofile", "r_emailaddress"}, config.Scopes, "scopes are not correct")
	assert.Equal(linkedin.Endpoint, config.Endpoint, "endpoint is not correct")
}

type HTTPClientMock struct{}

var i int

func (_m HTTPClientMock) Get(url string) (*http.Response, error) {
	// First call is successful.
	if i == 0 {
		i++
		return &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("Successful data"))),
		}, nil
	}

	// Second call is forbidden.
	if i == 1 {
		i++
		return &http.Response{
			Status:     "403 Forbidden",
			StatusCode: 403,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("Forbidden"))),
		}, nil
	}

	// Third call is bad request.
	if i == 2 {
		i++
		return &http.Response{
			Status:     "400 Bad Request",
			StatusCode: 400,
			Body:       ioutil.NopCloser(bytes.NewReader([]byte("Bad request"))),
		}, nil
	}

	// Following calls return an error.
	return nil, fmt.Errorf("Error on request")
}
