package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/handlers"
	"github.com/valorvig/bookings/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig        // put it here to let "middleware" which is in the same package "main" can also access to this "app"
var session *scs.SessionManager // customized session

func main() {
	// change this to true when in production
	app.InProduction = false // change one place but affect everywhere

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
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

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
