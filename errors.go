package catena

import (
	"fmt"
	"net/http"
)

// Errorf returns a new ErrorHandler with the specified code and message format.
func Errorf(status int, message string, a ...interface{}) error {
	return &ErrorHandler{
		status:  status,
		message: fmt.Sprintf(message, a...),
	}
}

// ErrorHandler implements http.Handler for writing JSON API errors to the client and
// also implements error so that it can be used as an error to return from handler
// methods and written to http responses in middleware.
type ErrorHandler struct {
	status  int    // http status to return with the application specific code
	message string // application specific error message
}

// Error implements error
func (e *ErrorHandler) Error() string {
	return fmt.Sprintf("[%d] %s", e.status, e.message)
}

// ServeHTTP replies to a request by writing the error message as json with the code
// and the message, writting the http status code. It does not otherwise end the request
// the caller should ensure no further writes are done to w.
func (e *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ctjson)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if e.status > 0 {
		w.WriteHeader(e.status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, `{"code": %d, "message": "%s"}`, e.status, e.message)
}

// Default ErrorHandlers for standard http request errors
var (
	NotFound         = &ErrorHandler{status: http.StatusNotFound, message: http.StatusText(http.StatusNotFound)}
	MethodNotAllowed = &ErrorHandler{status: http.StatusMethodNotAllowed, message: http.StatusText(http.StatusMethodNotAllowed)}
)

// PanicHandler allows the application to recover from panics and
func PanicHandler(w http.ResponseWriter, r *http.Request, ctx interface{}) {
	// TODO: add Sentry integration here
	// Add panic logging and recovering here
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
