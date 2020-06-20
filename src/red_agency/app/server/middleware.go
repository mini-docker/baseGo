package server

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c Context) bool
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(Context) bool {
	return false
}
