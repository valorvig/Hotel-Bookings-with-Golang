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

// create a connection pool for my sql to postgres db
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
