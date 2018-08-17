package mocks

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// OAuth2Config is the mock structure for linkedin.OAuth2Config.
type OAuth2Config struct{}

var token *oauth2.Token

// Exchange mocks the exchange of a code for a token.
func (_m OAuth2Config) Exchange(ctx context.Context, code string,
	opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {

	switch code {
	case "correct_code_succesful_data_retrieval":
		token = &oauth2.Token{AccessToken: "token_enable_data_retrieval"}
		return token, nil

	case "correct_code_failed_data_retrieval":
		token = &oauth2.Token{AccessToken: "token_disable_data_retrieval"}
		return token, nil

	default:
		return nil, fmt.Errorf("Invalid code")
	}

}

// AuthCodeURL mocks the generation of the authentication URL.
func (_m OAuth2Config) AuthCodeURL(state string,
	opts ...oauth2.AuthCodeOption) string {

	return ""
}

// Client mocks the return of http.Client
func (_m OAuth2Config) Client(ctx context.Context,
	t *oauth2.Token) *http.Client {

	return &http.Client{}
}

// GetAccessToken returns the access token value.
func GetAccessToken() string {
	return token.AccessToken
}
