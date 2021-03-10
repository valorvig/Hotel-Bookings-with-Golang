package dbrepo

import (
	"database/sql"

	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/repository"
)

// below is an example of how to conenct to postgres db

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB // holds a database connection pool
}

type testDBRepo struct { // has same fields as postgresDBRepo
	App *config.AppConfig
	// we aren't going to populate with anything, but it needs to exist. So we can call database funcitons that don't actually have a database behind them.
	DB *sql.DB
}

// create a connection pool for my sql to postgres db
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingsRepo(a *config.AppConfig) repository.DatabaseRepo {
	// need to create test-repo.go and copy all the db methods from postgres.go to there, otherwise &testDBRepo will be marked error.
	return &testDBRepo{
		App: a,
	}
}
