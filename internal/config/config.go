package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/valorvig/bookings/internal/models"
)

/*
handlers want to access "session" in the package main,
but we can't import main to handlers since it will cause import cycle.
So this package config is like a medium for such a case.
*/

// AppConfig holds the application config
// declare here to be available for everyone that has access to AppConfig
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger // write and read the log file to anywhere - terminal, Windows, database
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData // we put MailData in the cofig so that it's available for the entire application
}
