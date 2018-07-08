package linkedin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gvso/cardenal/constants"
	"github.com/gvso/cardenal/settings"

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

var conf *oauth2.Config

// Login starts authorization request to LinkedIn.
func Login(w http.ResponseWriter, r *http.Request) {
	conf = getConfig()

	fmt.Println(conf.RedirectURL)
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, 303)
}

// Callback handles LinkedIn API callback.
func Callback(w http.ResponseWriter, r *http.Request) {
	conf = getConfig()

	ctx := context.Background()

	// If there was an error when authenticating.
	if r.URL.Query()["error"] != nil {
		// @TODO: redirect to error page.
		fmt.Fprint(w, "Could not login.")

		return
	}

	code := r.URL.Query()["code"][0]

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Println(err)

		// @TODO: redirect to error page.
		fmt.Fprint(w, "Could not login.")

		return
	}

	client := conf.Client(ctx, tok)

	data, err := getProfile(client)

	// getData returns an error if request to API was not successful.
	if err != nil {
		if settings.Development {
			log.Println(err)
		}

		// @TODO: redirect to error page.
		fmt.Fprint(w, "Could not login.")

		return
	}

	w.Header().Add("Content-Type", "application/json")

	// At this point, we know data is a slice of bytes. Convert it to string.
	d := string(data)

	dataMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(d), &dataMap)
	if err != nil {
		panic(err)
	}

	//maputils.DumpMap("", dataMap)

	fmt.Fprint(w, d)
}

// Gets profile data from LinkedIn.
func getProfile(client HTTPClient) ([]byte, error) {
	f := strings.Join(fields, ",")

	// Request data from API.
	resp, err := client.Get(constants.LinkedInBaseURL + "/people/~:(" + f + ")?format=json")

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

// Returns the OAuth2 client configuration.
func getConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     settings.LinkedIn.ClientID,
		ClientSecret: settings.LinkedIn.ClientSecret,
		RedirectURL:  settings.LinkedIn.RedirectURLHost + ":" + settings.Port + constants.LinkedInRedirectPath,
		Scopes: []string{
			"r_basicprofile",
			"r_emailaddress",
		},
		Endpoint: linkedin.Endpoint,
	}
}

// HTTPClient is an interface for HTTP clients.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}
