package mocks

// GinContext is the mock structure for GinContext.
type GinContext struct {
	Token string
}

var redirectCalled = false
var setCookieCalled = false

// Redirect mocks an HTTP redirection call.
func (_m GinContext) Redirect(code int, location string) {
	redirectCalled = true
}

// SetCookie mocks a SetCookie call.
func (_m *GinContext) SetCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) {
	setCookieCalled = true

	_m.Token = value
}

// WasRedirectedCalled determines if Redirect was called.
func (_m GinContext) WasRedirectedCalled() bool {
	return redirectCalled
}

// WasSetCookieCalled determines if Redirect was called.
func (_m GinContext) WasSetCookieCalled() bool {
	return setCookieCalled
}
