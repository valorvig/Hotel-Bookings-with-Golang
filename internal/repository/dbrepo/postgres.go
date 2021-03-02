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
// After adding "returning id", postgres now returns not just an error, but both an int (id) and error (err)
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
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

	var newID int

	stmt := `insert into reservations (first_name, last_name, email, phone, start_date,
		end_date, room_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	// psotgres uses "returning id" instead of "last_insert_id"

	// With context, much safer and more robust means of talking to the database
	// _, err := m.DB.ExecContext(ctx, stmt,
	err := m.DB.QueryRowContext(ctx, stmt, // QueryRowContext returns *Row
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(), // the current time we have created this
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *postgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id,
		created_at, updated_at, restriction_id)
		values 
		($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		time.Now(),
		time.Now(),
		r.RestrictionID,
	)

	if err != nil {
		return err
	}
	return nil
}
