package globalmocks

import (
	"fmt"
)

// GinContext is the mock structure for GinContext.
type GinContext struct {
	cookieValues map[string]string

	AbortCall     abort
	CookieCall    cookie
	NextCall      next
	RedirectCall  redirect
	SetCall       set
	SetCookieCall setCookie
	StringCall    stringCall
}

// Abort mocks a call to Abort
func (_m *GinContext) Abort() {
	_m.AbortCall = abort{true}
}

// Cookie mocks a call to Cookie
func (_m *GinContext) Cookie(name string) (string, error) {
	_m.CookieCall = cookie{true, name}

	if value, ok := _m.cookieValues[name]; ok {
		return value, nil
	}

	return "", fmt.Errorf("No cookie with such value exists")
}

// Next mocks a call to Next
func (_m *GinContext) Next() {
	_m.NextCall = next{true}
}

// Redirect mocks a call to Redirect.
func (_m *GinContext) Redirect(code int, location string) {
	_m.RedirectCall = redirect{true, code, location}
}

// Set mocks a call to Set
func (_m *GinContext) Set(key string, value interface{}) {
	_m.SetCall = set{true, key, value}
}

// SetCookie mocks a call to SetCookie.
func (_m *GinContext) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	_m.SetCookieCall = setCookie{true, name, value, maxAge, path, domain, secure, httpOnly}

	_m.cookieValues = make(map[string]string)
	_m.cookieValues[name] = value
}

// String mocks a call to String.
func (_m *GinContext) String(code int, format string, values ...interface{}) {
	_m.StringCall = stringCall{true, code, format, values}
}

/*******************************************************************************
** Defines structs to check if functions were called with expected parameters **
*******************************************************************************/
type abort struct {
	Called bool
}
type cookie struct {
	Called bool
	Name   string
}
type next struct {
	Called bool
}
type redirect struct {
	Called   bool
	Code     int
	Location string
}
type set struct {
	Called bool
	Key    string
	Value  interface{}
}
type setCookie struct {
	Called   bool
	Name     string
	Value    string
	MaxAge   int
	Path     string
	Domain   string
	Secure   bool
	HTTPOnly bool
}
type stringCall struct {
	Called bool
	Code   int
	Format string
	Values []interface{}
}
