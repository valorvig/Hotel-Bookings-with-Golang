// connect our application to the database

package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn" // gives us access to the subroutines
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB holds the database connection pool
// by using struct, it can hold a driver to postgres or whatever through SQL field
type DB struct {
	SQL *sql.DB // database connection
}

var dbConn = &DB{} // reference to the DB type and initialize and empty struct

// Define the nature of our connection pool
// Maximum number of open database connection I could have
const maxOpenDbConn = 10 // not allow to ahve more than 10 in a given time is reasonable in this case
// How many connection can remain in the pool but remain idle
const maxIdleDbConn = 5 // 5 is nice
// Maximum lifetime for a database connection
const maxDbLifetime = 5 * time.Minute

// ConnectSQL creates database pool for Postgres
func ConnectSQL(dsn string) (*DB, error) { // dsn - data source name
	d, err := NewDatabase(dsn)
	if err != nil {
		// let it die beacause we can't go any further with our application if it can't connect to the database
		panic(err) // panic - if this doesn't work, I don't want anything to start up. THe program just doesn't work.
	}
	d.SetMaxOpenConns(maxOpenDbConn)    // stop it from growing out of control
	d.SetMaxIdleConns(maxIdleDbConn)    // remove idle conenctions and return the db when they're not being used
	d.SetConnMaxLifetime(maxDbLifetime) // ensure they have a certain lifetime for all of our db conenctions

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// testDB tries to ping database
func testDB(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("gpx", dsn)
	if err != nil {
		return nil, err
	}

	// check all in one step/line
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
