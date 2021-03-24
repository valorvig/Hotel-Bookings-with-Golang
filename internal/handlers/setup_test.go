// set up our testing environment

package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/models"
	"github.com/valorvig/bookings/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"
var functions = template.FuncMap{ // var functions = template.FuncMap{} // put this as the same one as from render.go to avoid error
	"humanDate":  render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate":    render.Iterate,
	"add":        render.Add,
}

// we won't just use getRoutes, but we'll take advantage of the test_main that we used in our other tests
// TestMain is part of the testing package in the standarad library
func TestMain(m *testing.M) {
	// Cut from getRoutes (that is laso from main.go)

	// register the models.Reservation to the session, so we know that we can use that in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// change this to true when in production
	app.InProduction = false

	// paste these from main.go to avoid handler error while testing
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // Ldate & Ltime are the local date and local time
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session // store session into the field

	// actual channel (the same as in main.go)
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	// this will duplicate the functionality that we do in our live application
	// we don't want to actually send the mail in our test
	listenForMail()

	tc, err := CreateTestTemplateCache() // we don't want to call CreateTemplateCache directly
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	// if false, it's going to rebuild the page (createTemplateCache) on every request in render
	// if true, then get the tempalte from the template cache
	// the problem is, if createTemplateCache is called, the pathToTemplates it used will be "./templates" from render.go, not the right ome for this test
	app.UseCache = true

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	render.NewRenderer(&app)

	//-----------------------
	// add os.Exit to run the test before it dies (exits)
	os.Exit(m.Run())
}

// doing pretty much the same that we do in our live application, but we want to skip the actual seding of mail
func listenForMail() {
	go func() {
		for {
			// handle but do nothing with anything we send to the mail chnannel
			_ = <-app.MailChan
		}
	}()
}

// This time we create a function that creates everything we need rather than using "TestMain" like we did in cmd/web.
// require routes to call the handlers
// routes return a Handler
func getRoutes() http.Handler {
	// copy from main.go ------------------------------------------------

	/*
		gob.Register(models.Reservation{}) // register the models.Reservation to the session, so we know that we can use that in the session

		// change this to true when in production
		app.InProduction = false

		// paste these from main.go to avoid handler error while testing
		infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // Ldate & Ltime are the local date and local time
		app.InfoLog = infoLog

		errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
		app.ErrorLog = errorLog

		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction

		app.Session = session // store session into the field

		tc, err := CreateTestTemplateCache() // we don't want to call CreateTemplateCache directly
		if err != nil {
			log.Fatal("cannot create template cache")
		}

		app.TemplateCache = tc
		// if false, it's going to rebuild the page (createTemplateCache) on every request in render
		// if true, then get the tempalte from the template cache
		// the problem is, if createTemplateCache is called, the pathToTemplates it used will be "./templates" from render.go, not the right ome for this test
		app.UseCache = true

		repo := NewTestRepo(&app)
		NewHandlers(repo)
		render.NewRenderer(&app)
	*/

	// copy from routes.go -----------------------------------------------

	mux := chi.NewRouter()

	// middleware stack will be firstly execute before searching for a matching route
	mux.Use(middleware.Recoverer) // recover from panic
	// mux.Use(WriteToConsole)

	// comment this out since we son't want to use the csrf token while we are testing with POST
	// mux.Use(NoSurf) // this will prevent all POST without passing the csrf protection

	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	// copy all the actual routes that don't exist herr yet and paste them here, then remove the "handlers." part

	mux.Get("/user/login", Repo.ShowLogin)
	mux.Post("/user/login", Repo.PostShowLogin)
	mux.Get("/user/logout", Repo.Logout)

	// don't need to be wrapped (in /admin routes) because we're just testing the handlers not the authentication
	// however, need to add the "/admin/" prefix for these routes
	mux.Get("/admin/dashboard", Repo.AdminDashboard)
	mux.Get("/admin/reservations-new", Repo.AdminNewReservations)
	mux.Get("/admin/reservations-all", Repo.AdminAllReservations)
	mux.Get("/admin/reservations-calendar", Repo.AdminReservationsCalendar)
	mux.Post("/admin/reservations-calendar", Repo.AdminPostReservationsCalendar)
	mux.Get("/admin/process-reservation/{src}/{id}/do", Repo.AdminProcessReservation) // the same problem occur so we need to add "do" like "show"
	mux.Get("/admin/delete-reservation/{src}/{id}/do", Repo.AdminDeleteReservation)
	mux.Get("/admin/reservations/{src}/{id}/show", Repo.AdminShowReservation)
	mux.Post("/admin/reservations/{src}/{id}", Repo.AdminPostShowReservation)

	// the img won't know how to get to the folder "/static/images/house.jpg"
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) // strip "/static" with "./static"

	return mux
}

// NoSurf adds CSRF protection to all POST request.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	// Create NoSurf token
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",                  // Cookie is available within this path
		Secure:   app.InProduction,     // false for development, true for production (only HTTPS allowed)
		SameSite: http.SameSiteLaxMode, // not strict to only the same site (with the same domain name)
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every request.
func SessionLoad(next http.Handler) http.Handler {
	// LoadAndSave provides middleware that automatically loads and saves session data for the current request, and communicates the session token to and from the client in a cookie
	return session.LoadAndSave(next)
}

// CreateTemplateCache creates a template chache as a map
// The purpose is for a structure to hold and look for things quickly
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{} // ready to use parsed templates

	// func Glob(pattern string) (matches []string, err error) - Glob returns the names of all files matching pattern or nil if there is no matching file.
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		// func Base(path string) string - Base returns the last element of path. Trailing path separators are removed before extracting the last element.
		name := filepath.Base(page) // name is about.page.tmpl or home.page.tmpl

		// fmt.Println("Page is currently", page) // Ex. Page is currently templates\about.page.tmpl

		ts, err := template.New(name).Funcs(functions).ParseFiles(page) // parse the current page into ts
		// fmt.Println("ts1: ", ts)
		// Home --> &{<nil> 0xc0000205c0 0xc00013c200 0xc00004e1e0}
		// About --> &{<nil> 0xc00008a180 0xc0000a6100 0xc0000ce000}
		if err != nil {
			// fmt.Println("ERROR: ts, err := template.New(name).Funcs(functions).ParseFiles(page)")
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		// fmt.Println("matches: ", matches) // [templates\base.layout.tmpl]
		if err != nil {
			// fmt.Println("ERROR:	matches, err := filepath.Glob(\"./templates/*.layout.tmpl\")")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates)) // parse the base layout to the current ts
			// fmt.Println("ts2: ", ts)
			// Home --> &{<nil> 0xc0000205c0 0xc00013c200 0xc00004e1e0}
			// About --> &{<nil> 0xc00008a180 0xc0000a6100 0xc0000ce000}
			if err != nil {
				// fmt.Println("ERROR: ts, err = ts.ParseGlob(\"./templates/*.layout.tmpl\")")
				return myCache, err
			}
		}

		// map the current ts (with both a specific page and the base layout) to its coordinated name
		myCache[name] = ts // [?] Map seems to copy value ts to new address
		// fmt.Println("myCache1: ", myCache)
	}
	// fmt.Println("myCache2: ", myCache) // map[about.page.tmpl:0xc000088390 home.page.tmpl:0xc000073530]
	return myCache, nil
}
