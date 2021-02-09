package main // same package "main" as main.go

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/valorvig/bookings/pkg/config"
	"github.com/valorvig/bookings/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// middleware stack will be firstly execute before searching for a matching route
	mux.Use(middleware.Recoverer) // recover from panic
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	// the img won't know how to get to the folder "/static/images/house.jpg"
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // strip "/static" with "./static"

	return mux
}
