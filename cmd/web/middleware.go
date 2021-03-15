package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/valorvig/bookings/internal/helpers"
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
		Path:     "/",                  // Cookie is available within this path
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

/*
In this middleware, I'm going to call the helper function IsAuthenticated,
and it requires a pointer to the http.request as a parameter,
so we need to have access to the request - return the http.handler
by calling the http.HandlerFunc and having as its parameter an anonymous function
that takes w and r
*/
// This middleware is something you can apply to routes that you want to protect and to ensure that only logged-in people have access to the routes.
// This is the same logic as the Nosurf and the SessionLoad, but it's our own custom middleware that actually has access to the request.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { // reutrn something that has access to ReponseWriter and Request from http.HandlerFunc
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r) // pass onto the next middleware if any
	})
}
