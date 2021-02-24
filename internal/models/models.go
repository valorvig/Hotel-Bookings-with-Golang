// store our models here that will be uesd to store things in a database, to get a data from a db and store them in these models

package models

// reservation holds reservation data
type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
