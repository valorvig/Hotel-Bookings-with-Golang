package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// WriteToConsole writes to console everytime a user accessing a page
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page")
		// move to the next middleware or the part of our file that returns our mux
		next.ServeHTTP(w, r)
	})
}

// NoSurf adds CSRF protection to all POST request.
// It creates token every time a user accessing a page.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// Create NoSurf token
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,     // false for development, true for production (only HTTPS allowed)
		SameSite: http.SameSiteLaxMode, // not strict to only the same site (with the same domain name)
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request.
// SessionLoad provides middleware from LoadAndSave (from scs/v2).
func SessionLoad(next http.Handler) http.Handler {
	// LoadAndSave provides middleware that automatically loads and saves session data for the current request, and communicates the session token to and from the client in a cookie
	return session.LoadAndSave(next)
}
