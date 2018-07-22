package mocks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GinContext is the mock structure for GinContext.
type GinContext struct {
	Request *http.Request
	Writer  gin.ResponseWriter
}

var redirectCalled = false

// Redirect mocks an HTTP redirection call.
func (_m GinContext) Redirect(code int, location string) {
	redirectCalled = true
}

func (_m GinContext) String(code int, format string, values ...interface{}) {

}

// WasRedirectedCalled determine if Redirect was called.
func (_m GinContext) WasRedirectedCalled() bool {
	return redirectCalled
}
