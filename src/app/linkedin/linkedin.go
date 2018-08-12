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
	"github.com/gvso/cardenal/src/app/jwt"
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
var fields = []string{"id", "first-name", "last-name", "headline", "industry",
	"picture-urls::(original)", "specialties", "positions", "public-profile-url"}

var conf OAuth2Config

var profileRetriever *ProfileRetriever

/**
 * Helper functions from external packages.
 *
 * For testing purposes, variables are declared which are defined as in external
 * packages. The advantage of doing this is that these variables can later be
 * overwritten for unit testing purposes.
 */
var processUserAuth = user.ProcessUserAuth

func init() {
	conf = getConfig().(*oauth2.Config)

	profileRetriever = NewProfileRetriever(getProfile)
}

// Login starts authorization request to LinkedIn.
//
// It generates the authentication url to LinkedIn and redirects the user for
// authorization.
func Login(c global.GinContext) {
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
func Callback(c *gin.Context) {
	ctx := context.Background()

	// If there was an error when authenticating.
	if c.Query("error") != "" {
		// @TODO: redirect to error page.
		c.String(http.StatusBadRequest, "Could not login.")

		return
	}

	code := c.Query("code")

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Println(err)

		// @TODO: redirect to error page.
		c.String(http.StatusBadRequest, "Could not login.")

		return
	}

	client := conf.Client(ctx, tok)

	data, err := profileRetriever.getProfile(client)
	if err != nil {
		if settings.Development {
			log.Println(err)
		}

		// @TODO: redirect to error page.
		c.String(http.StatusBadRequest, "Could not login.")

		return
	}

	processSuccessfulAuth(c, data)
}

// Manages workflow when authentication and data gathering have succeded.
//
// It process the user data and authentication workflow. Then, it sets or
// updates the token in cookies.
func processSuccessfulAuth(c global.GinContext, data []byte) {
	user, err := processUserAuth(data)

	err = setCookie(c, user)
	if err != nil {
		panic(err)
	}

	c.JSON(200, user)
}

// Gets profile data from LinkedIn.
//
// It is used by Callback to request basic user information from LinkedIn API.
func getProfile(client HTTPClient) ([]byte, error) {
	f := strings.Join(fields, ",")

	// Request data from API.
	resp, err := client.Get(global.LinkedInBaseURL + "/people/~:(" + f + ")?format=json")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {

		if settings.Development {
			log.Println(string(data))
		}

		return nil, fmt.Errorf("request was not successful")
	}

	return data, nil
}

// Sets token value in cookie.
//
// It gets the user firstname, lastanme and id on LinkedIn and create a JWT
// token which is later stored in cookie that expires in 7 days.
func setCookie(c global.GinContext, user map[string]string) error {
	token, err := jwt.CreateToken(user)

	if err != nil {
		return err
	}

	c.SetCookie("token", token, timeutils.GetSeconds(7), "/", "", false, true)

	return nil
}

// Returns the OAuth2 client configuration.
//
// It is used to generate a valide oauth.Config which is later used by Login and
// Callback to connect to LinkedIn.
func getConfig() OAuth2Config {
	return &oauth2.Config{
		ClientID:     settings.LinkedIn.ClientID,
		ClientSecret: settings.LinkedIn.ClientSecret,
		RedirectURL:  settings.LinkedIn.RedirectURLHost + ":" + settings.Port + global.LinkedInRedirectPath,
		Scopes: []string{
			"r_basicprofile",
			"r_emailaddress",
		},
		Endpoint: linkedin.Endpoint,
	}
}

type profileGetter func(client HTTPClient) ([]byte, error)

// ProfileRetriever implements methods to return LinkedIn profile data.
type ProfileRetriever struct {
	profileGetter
}

// NewProfileRetriever return a reference to a ProfileRetriever.
func NewProfileRetriever(pg profileGetter) *ProfileRetriever {
	return &ProfileRetriever{profileGetter: pg}
}

func (p *ProfileRetriever) getProfile(client HTTPClient) ([]byte, error) {
	return p.profileGetter(client)
}

// OAuth2Config is an interface for oauth2.Config
type OAuth2Config interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, t *oauth2.Token) *http.Client
}

// HTTPClient is an interface for HTTP clients.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}
