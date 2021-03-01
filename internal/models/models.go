// store our models here that will be uesd to store things in a database, to get a data from a db and store them in these models

package models

import "time"

// Reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// describe our database in a format that Go can understand
// table users

// Users is the user model
type Users struct {
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

// Rooms is the room model
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions is the restriction model
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Reservations is  the reservation model
type Reservations struct {
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
	// We don't just have to put in the fields exactly as they exist in the database table
	Room Rooms
}

// RoomRestrictions is the room restriction model
type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	ReservationID int
	RiestrictonID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Room          Rooms
	// What if I want to also show the reservation and the restriction, just add them here
	Reservation Reservations
	Restriction Restrictions
}
