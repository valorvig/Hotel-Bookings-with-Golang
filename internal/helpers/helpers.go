package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/valorvig/bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a // initialize app config from main.go
}

// when you have a client error --> then (1) have the response writer to writer something to the client and (2) status to show what kind of error that was, ex. 400 or 403
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) { // in case something went wrong with the server
	// [Big] Stack tracing is useful for diagnostic
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())                                     // (1) the nature of the error itself and (2) the actual stack trace - the detailed info about the nature of the rror that took palce
	app.ErrorLog.Println(trace)                                                                    // you may send the error via email or others instead
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // use StatusText to make it absolutely technically correct
}
