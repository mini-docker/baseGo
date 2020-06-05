package middleware

import (
	"baseGo/src/fecho/golog"
	"baseGo/src/red_admin/app/server"
	"fmt"
	"runtime"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper server.Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `json:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `json:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `json:"disable_print_stack"`
	}
)

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	DefaultRecoverConfig = RecoverConfig{
		Skipper:           server.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() server.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config RecoverConfig) server.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(c server.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						golog.Error("middleware", "RecoverWithConfig", "error:%v", fmt.Errorf("[%s] %s %s\n", "PANIC RECOVER", err, stack[:length]))
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
