package mocks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPClient is the mock structure for linkedin.HTTPClient.
type HTTPClient struct{}

var i int

// Get mocks an HTTP get request.
func (_m HTTPClient) Get(url string) (*http.Response, error) {
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
