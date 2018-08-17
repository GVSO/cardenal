package linkedin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwtlibrary "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"

	"github.com/gvso/cardenal/src/app/global"
	"github.com/gvso/cardenal/src/app/global/mocks"
	"github.com/gvso/cardenal/src/app/jwt"
	"github.com/gvso/cardenal/src/app/linkedin/mocks"
	"github.com/gvso/cardenal/src/app/settings"
)

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	context := &globalmocks.GinContext{}

	Login(context)

	redirectURL := "https://www.linkedin.com/oauth/v2/authorization?access_type=offline&client_id=&redirect_uri=%3A%2Fapi%2Fservices%2Flogin%2Fcallback&response_type=code&scope=r_basicprofile+r_emailaddress&state=state"

	// Asserts that function was called with correct arguments.
	assert.True(context.RedirectCall.Called)
	assert.Equal(302, context.RedirectCall.Code)
	assert.Equal(redirectURL, context.RedirectCall.Location)
}

func TestCallback(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	assert := assert.New(t)

	router := setupRouter()

	testErrorParam(assert, router)

	conf = mocks.OAuth2Config{}
	testCallbackWithWrongCode(assert, router)

	testCallbackWithCorrectCode(assert, router)
}

func TestProcessSuccessfulAuth(t *testing.T) {

	assert := assert.New(t)

	c := &globalmocks.GinContext{}

	// Saves current function and restores it at the end.
	old := processUserAuth
	defer func() { processUserAuth = old }()

	userMap := map[string]string{
		"first_name": "John",
		"last_name":  "Smith",
		"id":         "JohnId123",
	}

	// Overwrites processUserAuth function.
	processUserAuth = func(user []byte) (map[string]string, error) {
		return userMap, nil
	}

	arg := []byte("{\"firstName\":\"John\",\"id\":\"JohnId123\",\"lastName\":\"Smith\"}")

	user, err := processSuccessfulAuth(c, arg)

	// Asserts the returning values of the function.
	assert.Nil(err)
	assert.Equal(userMap, user)

	// Asserts the SetCookie was called correctly.
	assert.True(c.SetCookieCall.Called)
	assert.Equal("token", c.SetCookieCall.Name)

	/////////////////////////////////////////////////
	// Tests when processUserAuth returns an error.//
	/////////////////////////////////////////////////
	oldProcessUserAuth := processUserAuth

	// Overwrites setCookie function.
	processUserAuth = func(user []byte) (map[string]string, error) {
		return nil, fmt.Errorf("error when processing user")
	}

	user, err = processSuccessfulAuth(c, arg)

	// Asserts the returning values of the function.
	assert.Nil(user)
	assert.Equal("error when processing user", err.Error())

	processUserAuth = oldProcessUserAuth

	//////////////////////////////////////////
	// Tests when setCokie returns an error.//
	//////////////////////////////////////////
	oldSetCookie := setCookie
	defer func() { setCookie = oldSetCookie }()

	// Overwrites setCookie function.
	setCookie = setCookieErrorMock
	fmt.Println("here")
	user, err = processSuccessfulAuth(c, arg)

	// Asserts the returning values of the function.
	assert.Nil(user)
	assert.NotNil(err)
}

func TestGetProfile(t *testing.T) {
	assert := assert.New(t)

	client := &mocks.HTTPClient{}

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

func TestSetCookie(t *testing.T) {
	assert := assert.New(t)

	context := &globalmocks.GinContext{}

	user := map[string]string{
		"firstName": "John",
		"lastName":  "Smith",
		"id":        "JohnId123",
	}

	err := setCookie(context, user)

	// Asserts that function was called with correct arguments.
	assert.Nil(err)
	assert.True(context.SetCookieCall.Called)
	assert.Equal("token", context.SetCookieCall.Name)
	assert.True(isTokenValid(context.SetCookieCall.Value))
	assert.Equal(604800, context.SetCookieCall.MaxAge)
	assert.Equal("/", context.SetCookieCall.Path)
	assert.Equal("/", context.SetCookieCall.Path)
	assert.False(context.SetCookieCall.Secure)
	assert.True(context.SetCookieCall.HTTPOnly)

	// Tests when createToken returns an error

	// Saves current function and restores it at the end.
	old := createToken
	defer func() { createToken = old }()

	// Overwrites createToken function.
	createToken = func(user map[string]string) (string, error) {
		return "", fmt.Errorf("could not create token")
	}

	err = setCookie(context, user)

	assert.Equal("could not create token", err.Error())
}
func TestGetConfig(t *testing.T) {

	// Overrides settings' values.
	settings.LinkedIn = settings.LinkedInConfig{
		ClientID:        "client123",
		ClientSecret:    "secret123",
		RedirectURLHost: "http://localhost",
	}
	settings.Port = "8000"

	config := getConfig().(*oauth2.Config)

	assert := assert.New(t)

	assert.Equal(settings.LinkedIn.ClientID, config.ClientID, "client ID is not correct")
	assert.Equal(settings.LinkedIn.ClientSecret, config.ClientSecret, "client secret is not correct")

	redirectURL := settings.LinkedIn.RedirectURLHost + ":" + settings.Port + global.LinkedInRedirectPath
	assert.Equal(redirectURL, config.RedirectURL, "redirect URL is not correct")

	assert.Equal([]string{"r_basicprofile", "r_emailaddress"}, config.Scopes, "scopes are not correct")
	assert.Equal(linkedin.Endpoint, config.Endpoint, "endpoint is not correct")
}

// Tests when LinkedIn returns the 'error' parameter.
//
// For instance, when user denied authorization, LinkedIn returns an error
// parameter, so user shouldn't be logged in in this case.
func testErrorParam(assert *assert.Assertions, router *gin.Engine) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/services/login/callback", nil)

	q := req.URL.Query()
	q.Add("error", "access_denied")
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("Could not login.", w.Body.String())
}

// Tests when the returned code from LinkedIn is not valid and token could not be
// generated
func testCallbackWithWrongCode(assert *assert.Assertions, router *gin.Engine) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/services/login/callback", nil)

	q := req.URL.Query()
	q.Add("code", "incorrect_code123")
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("Could not login.", w.Body.String())
}

// Tests when the returned code from LinkedIn is valid and token could be
// generated
func testCallbackWithCorrectCode(assert *assert.Assertions, router *gin.Engine) {
	// Saves current function and restores it at the end.
	old := getProfile
	defer func() { getProfile = old }()

	// Overwrites getProfile function.
	getProfile = getProfileMock

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/services/login/callback", nil)

	testSuccessfulDataRetrieval(assert, router, w, req)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/services/login/callback", nil)

	testFailedDataRetrieval(assert, router, w, req)
}

// Tests when data could be retrieved from LinkedIn.
func testSuccessfulDataRetrieval(assert *assert.Assertions, router *gin.Engine, w *httptest.ResponseRecorder, req *http.Request) {
	q := req.URL.Query()
	q.Add("code", "correct_code_succesful_data_retrieval")
	req.URL.RawQuery = q.Encode()

	// Saves current function and restores it at the end.
	old := processUserAuth
	defer func() { processUserAuth = old }()

	userMap := map[string]string{
		"first_name": "John",
		"last_name":  "Smith",
		"id":         "JohnId123",
	}

	// Overwrites processUserAuth function.
	processUserAuth = func(user []byte) (map[string]string, error) {
		return userMap, nil
	}

	router.ServeHTTP(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("{\"first_name\":\"John\",\"id\":\"JohnId123\",\"last_name\":\"Smith\"}", w.Body.String())

	////////////////////////////////////////////////////////
	// Tests when processSuccessfulAuth returns an error. //
	////////////////////////////////////////////////////////

	// Saves current function and restores it at the end.
	oldProcessSuccessfulAuth := processSuccessfulAuth
	defer func() { processSuccessfulAuth = oldProcessSuccessfulAuth }()

	// Overwrites processUserAuth function.
	processSuccessfulAuth = func(c global.GinContext, data []byte) (map[string]string, error) {
		return nil, fmt.Errorf("could not set cookie")
	}

	// Creates new recorder
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("Could not login.", w.Body.String())
}

// Tests when data could not be retrieved from LinkedIn even if token was
// succesfully generated.
func testFailedDataRetrieval(assert *assert.Assertions, router *gin.Engine, w *httptest.ResponseRecorder, req *http.Request) {
	settings.Development = true

	q := req.URL.Query()
	q.Add("code", "correct_code_failed_data_retrieval")
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("Could not login.", w.Body.String())
}

// Checks that generated token stored in cookie is valid.
//
// It parses the token string and check for errors and validity of token. If no
// errors are produced and token is valid, return true.
func isTokenValid(tokenString string) bool {
	token, err := jwtlibrary.Parse(tokenString, jwt.KeyFunction)

	if err != nil {
		return false
	}

	if token.Valid {
		return true
	}

	return false
}

// Sets up router to test callback.
func setupRouter() *gin.Engine {
	router := gin.Default()

	services := router.Group("/api/services")
	{
		services.GET("/login/callback", Callback)
	}

	return router
}

// Mocks getProfile function
func getProfileMock(client HTTPClient) ([]byte, error) {

	switch mocks.GetAccessToken() {
	case "token_enable_data_retrieval":
		return []byte("{\"firstName\":\"John\", \"id\":\"JohnId123\", \"lastName\": \"Smith\"}"), nil

	case "token_disable_data_retrieval":
		return nil, fmt.Errorf("data could not be retrieved")

	default:
		return nil, fmt.Errorf("unexpected error")
	}
}

// Mocks setCookie that returns an error.
func setCookieErrorMock(c global.GinContext, user map[string]string) error {
	return fmt.Errorf("could not set cookie")
}
