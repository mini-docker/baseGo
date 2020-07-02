package middleware

import (
	"baseGo/src/red_robot/app/server"
)

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c server.Context) bool
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(server.Context) bool {
	return false
}