package global

// GinContext is an interface for Gin Framework Context.
type GinContext interface {
	Abort()
	Cookie(name string) (string, error)
	Get(key string) (value interface{}, exists bool)
	JSON(code int, obj interface{})
	Next()
	Redirect(code int, location string)
	Set(key string, value interface{})
	SetCookie(name string, value string, maxAge int, path string, domain string,
		secure bool, httpOnly bool)
	String(code int, format string, values ...interface{})
}
