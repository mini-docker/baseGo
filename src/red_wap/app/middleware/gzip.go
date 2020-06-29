package middleware

import (
	"baseGo/src/red_wap/app/server"
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type (
	// GzipConfig defines the config for Gzip middleware.
	GzipConfig struct {
		// Gzip compression level.
		// Optional. Default value -1.
		Level   int `yaml:"level"`
		Skipper Skipper
	}

	gzipResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

const (
	gzipScheme = "gzip"
)

var (
	// DefaultGzipConfig is the default Gzip middleware config.
	DefaultGzipConfig = GzipConfig{
		Level:   -1,
		Skipper: DefaultSkipper,
	}
)

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip() server.MiddlewareFunc {
	return GzipWithConfig(DefaultGzipConfig)
}

// GzipWithConfig return Gzip middleware with config.
// See: `Gzip()`.
func GzipWithConfig(config GzipConfig) server.MiddlewareFunc {
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}

	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(c server.Context) error {

			if config.Skipper(c) {
				return next(c)
			}

			res := c.Response()
			res.Header().Add(server.HeaderVary, server.HeaderAcceptEncoding)
			if strings.Contains(c.Request().Header.Get(server.HeaderAcceptEncoding), gzipScheme) {
				res.Header().Set(server.HeaderContentEncoding, gzipScheme) // Issue #806
				rw := res.Writer
				w, err := gzip.NewWriterLevel(rw, config.Level)
				if err != nil {
					return err
				}
				defer func() {
					if res.Size == 0 {
						if res.Header().Get(server.HeaderContentEncoding) == gzipScheme {
							res.Header().Del(server.HeaderContentEncoding)
						}
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						res.Writer = rw
						w.Reset(ioutil.Discard)
					}
					w.Close()
				}()
				grw := &gzipResponseWriter{Writer: w, ResponseWriter: rw}
				res.Writer = grw
			}
			return next(c)
		}
	}
}

func (w *gzipResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent { // Issue #489
		w.ResponseWriter.Header().Del(server.HeaderContentEncoding)
	}
	w.Header().Del(server.HeaderContentLength) // Issue #444
	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get(server.HeaderContentType) == "" {
		w.Header().Set(server.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	w.Writer.(*gzip.Writer).Flush()
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *gzipResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
