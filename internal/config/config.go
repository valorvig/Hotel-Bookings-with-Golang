package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

/*
handlers want to access "session" in the package main,
but we can't import main to handlers since it will cause import cyvle.
So this package config is like a medium for such a case.
*/

// AppConfig holds the applicaiton config
// declare here to be available for everyone that has access to AppConfig
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger // write and read the log file to anywhere - terminal, Windows, database
	InProduction  bool
	Session       *scs.SessionManager
}
