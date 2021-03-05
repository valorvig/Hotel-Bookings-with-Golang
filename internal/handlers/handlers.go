package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/driver"
	"github.com/valorvig/bookings/internal/forms"
	"github.com/valorvig/bookings/internal/helpers"
	"github.com/valorvig/bookings/internal/models"
	"github.com/valorvig/bookings/internal/render"
	"github.com/valorvig/bookings/internal/repository"
	"github.com/valorvig/bookings/internal/repository/dbrepo"
)

// we can't have the struct model templatedata here since it's gpnna create import cycle error

// Repo the repository used by the handlers - it's implemented in routes.go
var Repo *Repository

// Repository is the repository type (Repository pattern)
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
// With this repository receiver, now this handler can access everything inside the repository, especially the AppConfig
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// grab the IP address of the person visiting the site and store it in the home page session
	// remoteIP := r.RemoteAddr // IPv4 or IPv6 address
	// m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	// m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	/*
		// example of using AllUsers
		// DB is a field of type "DatabaseRepo" which has method "AllUsers()"
		m.DB.AllUsers()
	*/

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hello, again."

	// // hold user's IP address in the session
	// // the value is empty if there is nothing in the session named "remote_ip"
	// remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	// after accessing the "session" from config, can do anything from it
	// m.App.Session.

	// stringMap["remote_ip"] = remoteIP // Ex. [::1]:18107 - "::1" is the loopback address in ipv6, equal to 127.0.0.1 in ipv4.

	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		// StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays form
// render the make-reservation tempalte and include the empty form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// var emptyReservation models.Reservation

	// Pull the ereservation out of the session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) // res - reservation
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID) // return the whole model of that room
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res) // put in session - used by PostReservation

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	// We have a stringMap in our models "TemplateData"
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	// data["reservation"] = emptyReservation // have to have the exact same name "reservation" as below
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil), // have access to the form the first time this page is loaded
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("can't get from session"))
		return
	}

	// it's a good practice to use ParseForm after parsing a form
	err := r.ParseForm()
	// err = errors.New("this is an error message") // intentionally create an error for testing purpose
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w, err)
		return
	}

	/*
		// We can't put these dates (string type) directly to StartDate/EndDate in models.Reservation because they are defined to accept type time
		sd := r.Form.Get("start_date")
		ed := r.Form.Get("end_date")

		// Format: 2020-01-01 --- 01/02 03:04:05PM '06 -0700
		// https://www.pauladamsmith.com/blog/2011/05/go_time.html
		// convert string to time
		layout := "2006-01-02"
		startDate, err := time.Parse(layout, sd)
		// fmt.Println("startDate====", reflect.TypeOf(startDate), startDate)
		if err != nil {
			helpers.ServerError(w, err)
		}
		endDate, err := time.Parse(layout, ed)
		// fmt.Println("endDate====", reflect.TypeOf(endDate), endDate)
		if err != nil {
			helpers.ServerError(w, err)
			return // (1) return to stop execution
		}

		// Convert string to int for RoomID in Reservation - Atoi (Alpha to integer)
		roomID, err := strconv.Atoi(r.Form.Get("room_id"))
		if err != nil {
			helpers.ServerError(w, err)
			return // (2) return to stop execution
		}
	*/

	// Update res
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	/*
		// We have to prevent users from losing all filled info after getting an error
		// So we need to indicate the error and where to fix it to the users
		// create "reservation" to reserve user's input data and prevent them from losing afterwards
		reservation := models.Reservation{
			FirstName: r.Form.Get("first_name"),
			LastName:  r.Form.Get("last_name"),
			Phone:     r.Form.Get("phone"),
			Email:     r.Form.Get("email"),
			StartDate: startDate,
			EndDate:   endDate,
			RoomID:    roomID,
		}
	*/

	// create a form with value
	form := forms.New(r.PostForm)

	// check if this form has first_name value, and add to the error mapping if it does
	// form.Has("first_name", r) // no longer need this after we've created "func (f *Form) Required(fields ...string)"
	form.Required("first_name", "last_name", "email") // form will have an error if one of these have an empty string
	form.MinLength("first_name", 3)
	// Required and then MinLength, so the first error from REquired will be displayed first
	form.IsEmail("email")

	// Form validation
	if !form.Valid() {
		data := make(map[string]interface{}) // create a variable to hold the data
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// After form validation, write the info to the database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return // (3) ***without return, it will still run, but won't stop execution at this point, and it will fail when it reach to the database insert part
	}

	// prepare for inserting restriction
	restriction := models.RoomRestriction{
		// ID:            0, // we don't need this - it will be automatically generated in the DB
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1, // have to fill up the db field before we can actually use the real thing
		// CreatedAt:     0, // these two will be taken care at the database level
		// UpdatedAt:     0,
		// Room:          0, // these three are not actually part of the db field, but we might use them in the future
		// Reservation:   0,
		// Restriction:   0,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// put reservation into session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	// shouldn't re-display the page - use redirect instead to prevent submitting the form twice accidentally
	// standard practice - anytime you receive a POST request, you should direct users to another page with an HTTP redirect
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther) // 303
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// After submitting and reload the page, the page is directed to this one
	// (In case you don't want the page to relaodm use another approach to let it run in the background)

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	// start := r.FormValue("start")
	// end := r.FormValue("end")

	// convert string to time
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return // (1) return to stop execution
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	/* for internal use and test only - telling something with a refresh page is not that userful
	// When there is availability, the actual rooms variable will ahve at least one entry
	for _, i := range rooms {
		m.App.InfoLog.Println("ROOM:", i.ID, i.RoomName)
	}
	*/

	if len(rooms) == 0 {
		// no availability
		// m.App.InfoLog.Println("No availability")
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return // want everything to stop at this point
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		// ID:        0, // don't care
		// FirstName: "",
		// LastName:  "",
		// Email:     "",
		// Phone:     "",
		StartDate: startDate,
		EndDate:   endDate,
		// RoomID:    0,
		// CreatedAt: time.Time{},
		// UpdatedAt: time.Time{},
		// Room:      models.Room{},
	}

	m.App.Session.Put(r.Context(), "reservation", res) // putting it in the session, now the info is available to use

	// w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// we don't need to use this JSON struct outside the package, and only its JSON fields is needed
type jsonResponse struct { // put it closr to the function you're using, so it's easy to find
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
// We're building a JSON request, not a web page that can't send back straight text
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, _ := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	resp := jsonResponse{
		OK:      available,
		Message: "",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w, err)
		return // let's return, we don't want to go any further
	}

	// log.Println(string(out))
	// Tell the browser, here's the kind of content (application JSON header) you're going to get
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact renders the search availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// get "reservation" out of the session
	// still not enough because the session, although it's storing a reservation, it has no idea what type that is - so we need to type assert
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// Ex. In case someone try to visit the page with "/reservation-summary" directly without making the reservation first
	if !ok {
		// 1) no session variable and data, 2) no variable called "reservation" in the session, 3) can't cast variable named "reservation" to a models.reservation
		// log.Println("cannot get item from session")
		m.App.ErrorLog.Println("Can't get error from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // use 300 because maybe they're going to make a reservation, or they'll come back later.
		return                                                 // we don't want to go further and display with a blank screen.
		// try typing "/reservation-summary" directly to see the result
	}

	// remove our data from the reservation
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id")) // get "id" from url parameter in routes page
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Pull the ereservation out of the session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) // res - reservation
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res) // put modified "res" back to the session

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther) // redirect to the Get page
}

/*
try using ./run.sh
*/
