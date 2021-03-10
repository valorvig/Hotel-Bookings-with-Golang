/*
After copy from postgres.go, paste it here.
Change all "postgresDBRepo" to "testDBRepo"

Back in our tests in dbrepo.goconst, we don't actuall passing the database connection pool

func NewTestingsRepo(a *config.AppConfig) repository.DatabaseRepo {
	// need to create test-repo.go and copy all the db methods from postgres.go to there.
	return &testDBRepo{
		App: a,
	}
}

Try leave every thing out to make it able to run first
*/

package dbrepo

import (
	"errors"
	"time"

	"github.com/valorvig/bookings/internal/models"
)

// Create and allow the methods to be used with only postgres repository at the moment

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
// After adding "returning id", postgres now returns not just an error, but both an int (id) and error (err)
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// for testing, we just don't want it to actually hot the database, so only return the newID

	/*
		// I can cancel context when something goes wrong
		// context.Background() is always available evertwhere in your application
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
	*/
	var newID int
	/*
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
		).Scan(&newID) // scan that row and put it into "newID"

		if err != nil {
			return 0, err
		}
	*/
	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	/*
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
	*/
	return nil // should be good for now
}

// SearchAvailabilityByDatesByRoomID returns true if availability exists for a given room (roomID), and false if no availability exists
// don't want to just give info about what rooms are available, which is only helpful for single room. Let's add "roomID"
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var numRows int

		query := `
			select
				count(id)
			from
				room_restrictions
			where
				room_id = $1
				and $2 < end_date and $3 > start_date;`

		row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
		err := row.Scan(&numRows)
		if err != nil {
			return false, err
		}

		if numRows == 0 {
			return true, nil
		}
	*/
	return false, nil
}

// SearchAvailabilityForAllRooms returns a slice of available rooms, if any, for given date range
// It would be more useful to return a slice of rooms rahter than just 2 integers for room_name and room_id.
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
	*/

	// only leave the declaration of model's room there
	var rooms []models.Room

	/*
		query := `
			select
				r.id, r.room_name
			from
				rooms r
			where r.id not in
			(select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);
			`
		// QueryContext returns roes, but QueryRowContext returns the lastest row
		rows, err := m.DB.QueryContext(ctx, query, start, end)
		if err != nil {
			return rooms, err
		}

		for rows.Next() {
			var room models.Room
			err := rows.Scan(
				&room.ID,
				&room.RoomName,
			)
			if err != nil {
				return rooms, err // we're going to ignore this rooms
			}

			rooms = append(rooms, room)
		}

		if err = rows.Err(); err != nil {
			return rooms, err
		}
	*/

	return rooms, nil
}

// GetRoomByID gets room by id
// We want to display the room name in summary page by creating a function GetRoomByID
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {
	/*
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
	*/

	var room models.Room

	// test error case
	if id > 2 {
		return room, errors.New("Some error")
	}

	/*
		query := `
			select id, room_name, created_At, updated_at from rooms where id = $1
			`
		// we only want one row
		row := m.DB.QueryRowContext(ctx, query, id)
		err := row.Scan(
			&room.ID,
			&room.RoomName,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return room, err
		}
	*/

	return room, nil
}
