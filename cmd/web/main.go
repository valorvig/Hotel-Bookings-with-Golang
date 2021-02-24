package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/handlers"
	"github.com/valorvig/bookings/internal/helpers"
	"github.com/valorvig/bookings/internal/models"
	"github.com/valorvig/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig        // put it here to let "middleware" which is in the same package "main" can also access to this "app"
var session *scs.SessionManager // customized session
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err) // stop the application from running
	}

	fmt.Println(fmt.Sprintf("starting application on port %s", portNumber))
	// _ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

/*
$go run cmd/web/*.go
*/

func run() error {
	// cut and paste everything in the main() starting and above "render.NewTemplates(&app)" till the beginning

	// what am I going to put in the session - the code with "gob" looks a bit odd
	// We can store primitives, but we need to actually tell our application about more complex types, structures
	// that we've defined ourselves and we want to store in the reservation model.
	// we've stored the value in the session. // m.App.Session.Put(r.Context(), "reservation", reservation)
	gob.Register(models.Reservation{}) // we want to store in the reservation model - now we can store values in the session

	// change this to true when in production
	app.InProduction = false // change one place but affect everywhere

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // Ldate & Ltime are the local date and local time
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// session := scs.New() // using short declaration means that the above session and this session are two different variables
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	// true - let cppkies persist after the browser is closed
	// false - the session will not persist the next time users open a window or fire up their web browser session cookie
	session.Cookie.Persist = true
	// what site this cookie applies to
	session.Cookie.SameSite = http.SameSiteLaxMode // use default - cookies were sent for all requests
	// cookie is encrypted and from only HTTPS
	session.Cookie.Secure = app.InProduction // false for development, truen for productiuon

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	// transfer main app and repo to other files
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	return nil
}
