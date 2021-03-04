// store our models here that will be uesd to store things in a database, to get a data from a db and store them in these models

package models

import "time"

// Describe our database in a format that Go can understand
// Table user from the database
// Use singular word since we only represent each stuff individually

// User is the user model
type User struct {
	// It's a convention to have field names of the user type correspond to the actual name used in the database.
	// Use uppercase letters for the first part of each name and leave the underscores out.
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restriction is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservation is  the reservation model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room // We don't just have to put in the fields exactly as they exist in the database table
}

// RoomRestriction is the room restriction model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Room
	// What if I want to also show the reservation and the restriction, just add them here
	Reservation Reservation
	Restriction Restriction
}
