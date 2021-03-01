package repository

import "github.com/valorvig/bookings/internal/models"

// make some methods or funcitons available to the repository for the database - NewRepo returns *Repository
type DatabaseRepo interface {
	AllUsers() bool

	// get info from the reservation page and put them to the database
	InsertReservation(res models.Reservation) error
}
