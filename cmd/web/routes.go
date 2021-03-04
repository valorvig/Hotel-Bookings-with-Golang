package main // same package "main" as main.go

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	// middleware stack will be firstly execute before searching for a matching route
	mux.Use(middleware.Recoverer) // recover from panic
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf) // this will prevent all POST without passing the csrf protection
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom) // whatever "/choose-room/..." it will be directed to ChooseRoom

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	// the img won't know how to get to the folder "/static/images/house.jpg"
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // strip "/static" with "./static"

	return mux
}
