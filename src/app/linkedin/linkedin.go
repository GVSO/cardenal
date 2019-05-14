package linkedin

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gvso/cardenal/src/app/global"
	"github.com/gvso/cardenal/src/app/settings"
	"github.com/gvso/cardenal/src/app/user"
	"github.com/gvso/cardenal/src/app/utils/timeutils"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

/**
 * User fields to be requested.
 *
 * @see https://developer.linkedin.com/docs/fields/basic-profile
 */
var fields = []string{"id", "localizedFirstName", "localizedLastName",
	"profilePicture(displayImage~:playableStreams)"}

var conf OAuth2Config

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as pointers to
 * function. The advantage of doing this is that these variables can later be
 * overwritten in the testing files.
 */
var processUserAuth = user.ProcessUserAuth

func init() {
	conf = getConfig().(*oauth2.Config)
}

// Login starts authorization request to LinkedIn.
//
// It generates the authentication url to LinkedIn and redirects the user for
// authorization.
var Login = func(c global.GinContext) {
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	c.Redirect(302, url)
}

// Callback handles LinkedIn API callback.
//
// After user authorizes (or not), user is redirected to our site where this
// function takes care of the authentication process.
//
// It checks if no error parameter has been provided by LinkedIn. Then, it
// exchanges the code for a token, which is later used to get user data from
// LinkedIn.
//
// When token and user data have been retrieved successfully, it calls
// processSuccessfulAuth to handle the authentication/registration workflow.
var Callback = func(c *gin.Context) {
	ctx := context.Background()

	// If there was an error when authenticating.
	if c.Query("error") != "" {
		processBadRequest(c, fmt.Errorf(c.Query("error")))

		return
	}

	// Gets token from LinkedIn.
	code := c.Query("code")

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		processBadRequest(c, err)

		return
	}

	// Gets user information from LinkedIn.
	client := conf.Client(ctx, tok)

	data, err := getProfile(client)
	if err != nil {
		processBadRequest(c, err)

		return
	}

	// Processes user authentication.
	userMap, err := processSuccessfulAuth(c, data, tok)
	if err != nil {
		processBadRequest(c, err)

		return
	}

	c.JSON(200, userMap)
}

var processBadRequest = func(c *gin.Context, err error) {
	if settings.Development {
		log.Println(err)
	}

	// @TODO: redirect to error page.
	c.String(http.StatusBadRequest, "Could not login.")
}

// Returns the OAuth2 client configuration.
//
// It is used to generate a valide oauth.Config which is later used by Login and
// Callback to connect to LinkedIn.
var getConfig = func() OAuth2Config {
	return &oauth2.Config{
		ClientID:     settings.LinkedIn.ClientID,
		ClientSecret: settings.LinkedIn.ClientSecret,
		RedirectURL:  settings.LinkedIn.RedirectURLHost + ":" + settings.Port + global.LinkedInRedirectPath,
		Scopes: []string{
			"r_liteprofile",
			"r_emailaddress",
		},
		Endpoint: linkedin.Endpoint,
	}
}

// Gets profile data from LinkedIn.
//
// It is used by Callback to request basic user information from LinkedIn API.
var getProfile = func(client HTTPClient) ([]byte, error) {
	projection := strings.Join(fields, ",")

	// Request data from API.
	resp, err := client.Get(global.LinkedInBaseURL + "/me/?projection=(" + projection + ")")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(string(data))
		if settings.Development {
			log.Println(string(data))
		}

		return nil, fmt.Errorf("request was not successful")
	}

	return data, nil
}

// Sets token value in cookie.
//
// It gets the token string and stores it in a cookie that expires in 7 days.
var setCookie = func(c global.GinContext, token string) {

	c.SetCookie("token", token, timeutils.GetSeconds(7), "/", "", false, true)
}

// Manages workflow when authentication and data gathering have succeded.
//
// It processes the user data and authentication workflow. Then, it sets or
// updates the token in cookie and database.
var processSuccessfulAuth = func(c global.GinContext,
	data []byte, linkedinToken *oauth2.Token) (map[string]string, error) {

	user, err := processUserAuth(data, linkedinToken)
	if err != nil {
		return nil, err
	}

	setCookie(c, user["token"])

	return user, nil
}
