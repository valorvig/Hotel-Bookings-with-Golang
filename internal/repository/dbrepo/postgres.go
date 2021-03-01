// put any functions that you want to be available to the interface (DatabaseRepo)

package dbrepo

import (
	"context"
	"time"

	"github.com/valorvig/bookings/internal/models"
)

// Create and allow the methods to be used with only postgres repository at the moment

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	/*
		Web applications are unpredictable in one sense:
		- The user might lose his or her connection to the internet partway through this transaction.
		- The user might just close their browser without warning in the midst of something that's taking place with the database in the background.
		So, we're going to use "context"
	*/
	// I can cancel context when something goes wrong
	// context.Background() is always available evertwhere in your application
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into reservation (first_name, last_name, email, phone, start date,
		end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	// With context, much safer and more robust means of talking to the database
	_, err := m.DB.ExecContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		time.Now, // the current time we have created this
		time.Now,
	)

	if err != nil {
		return err
	}

	return nil
}
