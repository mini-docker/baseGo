package middleware

import "github.com/mini-docker/baseGo/fecho/echo"

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c echo.Context) bool
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(echo.Context) bool {
	return false
}
