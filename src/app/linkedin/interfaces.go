package linkedin

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// OAuth2Config is an interface for oauth2.Config
type OAuth2Config interface {
	AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string,
		opts ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	Client(ctx context.Context, t *oauth2.Token) *http.Client
}

// HTTPClient is an interface for HTTP clients.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}
