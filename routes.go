package catena

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Routes creates and configures the server multiplexer with the API endpoints and
// methods, it does not add any additional middleware and primarily serves as
// documentation for how the API is configured.
func Routes() *httprouter.Router {
	// Create new httprouter with settings
	mux := httprouter.New()
	mux.RedirectTrailingSlash = true
	mux.RedirectFixedPath = true
	mux.HandleMethodNotAllowed = true
	mux.GlobalOPTIONS = nil

	// Handle routing errors and panics
	mux.NotFound = NotFound
	mux.MethodNotAllowed = MethodNotAllowed
	mux.PanicHandler = PanicHandler

	// Create basic routes
	mux.GET("/status/", status)
	return mux
}

func status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := make(map[string]interface{})
	status["status"] = "ok"
	status["timestamp"] = time.Now().Format(time.RFC3339Nano)
	status["version"] = Version

	data, err := json.Marshal(status)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", ctjson)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
