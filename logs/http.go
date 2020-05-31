package logs

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// NewHTTPLogger returns a web server specific logger that colorizes output based on
// status code rather than log level. This logger also wraps a handler and serves as
// both the request and response handling middleware.
func NewHTTPLogger(prefix string, handler http.Handler) *HTTPLogger {
	return &HTTPLogger{Logger: *New(prefix), handler: handler}
}

// HTTPLogger is a specialized logging module that colorizes output based on the status
// code of the response and uses levels slightly differently. HTTP status 100 levels are
// LevelInfo, status 200 or 300 are LevelStatus, and 400 or 500 are LevelWarn. You can
// use the same levels to surpress output if only errors are desirable to track. The
// primary difference between this logger with the default logger is that the
// colorization is not based on level but rather by status code. It also implements
// http.Handler so that it can be used in middlware.
type HTTPLogger struct {
	Logger
	handler http.Handler
}

func (l *HTTPLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle the incoming request
	ts := time.Now()            // track how much time it takes to serve the response
	lw := &responseLogger{w: w} // create a logging response writer to track size and status

	// Middleware is in the response path, so serve the next handler first.
	l.handler.ServeHTTP(lw, r)

	// Determine if log message is necessary
	status := lw.Status()
	if StatusLevel(status) < l.Logger.logLevel {
		return
	}

	// Log the request and the response.
	// [31/May/2020 08:11:06] "GET /api/status/ HTTP/1.1" 200 94
	// TODO: add tracing to the log output
	var buf strings.Builder
	estlen := len(l.Logger.prefix) + len(l.Logger.timestamp) + len(r.Method) + len(r.URL.Path) + len(r.Proto) + 8
	if l.Logger.colorize {
		estlen += 8
	}
	buf.Grow(estlen)

	// Colorize the output
	if l.Logger.colorize {
		buf.WriteString(statuscolors[status/100])
	}

	// Write the prefix
	buf.WriteString(l.Logger.prefix)

	// Write the timestamp
	if l.timestamp != "" {
		buf.WriteString(ts.Format(l.Logger.timestamp))
	}

	// Write the log message
	// TODO: better common logging format
	fmt.Fprintf(&buf, "\"%s %s %s\" %d %d", r.Method, r.URL.Path, r.Proto, status, lw.Size())

	// Reset the colorization
	if l.colorize {
		buf.WriteString(colorReset)
	}

	// Log the message
	l.Logger.logger.Print(buf.String())
}

// StatusLevel returns the log level for the specified http status code.
func StatusLevel(status int) uint8 {
	switch {
	case status < 200:
		return LevelInfo
	case status < 400:
		return LevelStatus
	case status < 600:
		return LevelWarn
	default:
		return LevelDebug
	}
}

// responseLogger wraps http.ResponseWriter to track status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}
