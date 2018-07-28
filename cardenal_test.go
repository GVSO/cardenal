package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServiceRoutes(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	testLoginRoute(t)
	testLoginCallbackRoute(t)
}

// Test the LinkedIn login route.
//
// We expect to get a redirection to LinkedIn service.
func testLoginRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/services/login", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 302, w.Code)
}

// Test the callback route after user has authenticated with LinkedIn.
//
// TODO: We might need to test all the cases as we do in the linkedin package.
func testLoginCallbackRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/services/login/callback", nil)
	router.ServeHTTP(w, req)

	// No correct code parameter was sent, so request should fail.
	assert.Equal(t, 400, w.Code)
}
