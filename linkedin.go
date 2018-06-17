package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/linkedin"
)

// LinkedInAPIException is the exception format LinkedIn API uses.
type LinkedInAPIException struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	RequestID string `json:"requestID"`
	Status    int    `json:"status"`
	Timestamp uint64 `json:"timestamp"`
}

func (e *LinkedInAPIException) String() string {
	return "{" +
		"\n\terrorCode: " + strconv.Itoa(e.ErrorCode) +
		"\n\tmessage: " + e.Message +
		"\n\trequestID " + e.RequestID +
		"\n\tstatus" + strconv.Itoa(e.Status) +
		"\n\ttimestamp" + strconv.FormatUint(e.Timestamp, 10) +
		"\n}"
}

// LinkedInException is a exception that ocurs when doing in LinkedIn requests.
type LinkedInException struct {
	Exception `json:"exception"`
	Response  LinkedInAPIException `json:"response"`
}

func (e *LinkedInException) Error() string {
	return e.Exception.Message + "\n" + e.Response.String()
}

const baseURL = "https://api.linkedin.com/v1"

/**
 * User fields to be requested.
 *
 * @see https://developer.linkedin.com/docs/fields/basic-profile
 */
var fields = []string{"id", "first-name", "last-name", "headline", "industry",
	"picture-urls::(original)", "specialties", "positions", "public-profile-url"}

var conf *oauth2.Config

func init() {
	conf = &oauth2.Config{
		ClientID:     Settings.LinkedIn.ClientID,
		ClientSecret: Settings.LinkedIn.ClientSecret,
		RedirectURL:  Settings.LinkedIn.RedirectURLHost + Settings.Port + "/api/services/login/callback",
		Scopes: []string{
			"r_basicprofile",
			"r_emailaddress",
		},
		Endpoint: linkedin.Endpoint,
	}
}

// Starts authorization request to LinkedIn.
func linkedinLogin(w *http.ResponseWriter, r *http.Request) {
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)

	http.Redirect(*w, r, url, 303)
}

// Handles LinkedIn API callback.
func linkedinCallback(w *http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	code := r.URL.Query()["code"][0]

	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Println(err)

		// @TODO: redirect to error page.
		return
	}

	client := conf.Client(ctx, tok)

	var exception LinkedInException
	data, err := getData(client, &exception)

	(*w).Header().Add("Content-Type", "application/json")

	// getData return an error if request to API was not successful.
	if err != nil {
		log.Println(err)

		// Prints in broser.
		if Settings.Development && exception.Response.Message != "" {
			json.NewEncoder(*w).Encode(&exception)
		}

		// @TODO: redirect to error page.
		return
	}

	// At this point, we know data is a slice of bytes. Convert it to string.
	d := string(data.([]byte))

	dataMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(d), &dataMap)
	if err != nil {
		panic(err)
	}

	//maputils.DumpMap("", dataMap)

	fmt.Fprint(*w, d)
}

// Gets data from LinkedIn.
func getData(client *http.Client, exception *LinkedInException) (interface{}, error) {
	f := strings.Join(fields, ",")

	// Request data from API.
	resp, err := client.Get(baseURL + "/people/~:(" + f + ")?format=json")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		// Populate LinkedInException fields.
		json.Unmarshal(data, &exception.Response)
		exception.Exception = Exception{Message: "Could not retrieved data from LinkedIn"}

		return nil, exception
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}
