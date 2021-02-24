package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/models"
)

// the firs thing I want to test is "func AddDefaultData" from "render.go." So we need *models.TemplateData and *http.Request as inputs.
// We aren't testing the entire application from main function, but only this one function at default data in isolation.
// That means we need to build this http.Request, and I can't just build an empty request.
// We need to build a request that has session data. wihout the session data, it will fail.

// create variable that render is going to need
var session *scs.SessionManager
var testApp config.AppConfig

// Try using the built-in function to test once again
// ***This function gets called before any of the tests are tun
func TestMain(m *testing.M) {
	// We need the session information in here. Copy from main.go and change "app." to our "testApp." instead.
	// What am I going to put in the session
	gob.Register(models.Reservation{})

	// Change this to true when in production
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime) // Ldate & Ltime are the local date and local time
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	// let's make sure that our testApp has things we need
	app = &testApp

	// --------------------------------------

	os.Exit(m.Run()) // run our test Run() just before it exits
}

// we are replicating our own struct representing ResponseWriter
/*
type ResponseWriter interface {
	Header() Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
*/

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	// staisfy the interface by creating an empty header variable
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {

}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b) // int can't be random but its length
	return length, nil
}
