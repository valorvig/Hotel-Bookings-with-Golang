package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/valorvig/bookings/internal/models"
)

// create a web server that returns status code and we can post to
// create a client that can call that server
// Golang will do all above for us
// we can just create the test server and store in a variable

// a variable to hold whatever we're posting to the page
type postData struct {
	key   string
	value string
}

// variable for the actual test
var theTests = []struct {
	name   string // name of the test
	url    string // the path which is matched by our roots
	method string // GET or POST
	// params             []postData // things that are being posted and that's a form can have more than 1 input
	expectedStatusCode int // what kind of response code we're getting back from the web server (200, 404, etc.)
}{
	/* GET */
	{"home", "/", "GET" /*[]postData{},*/, http.StatusOK},
	{"about", "/about", "GET" /*[]postData{},*/, http.StatusOK},
	{"gq", "/generals-quarters", "GET" /*[]postData{},*/, http.StatusOK},
	{"ms", "/majors-suite", "GET" /*[]postData{},*/, http.StatusOK},
	{"sa", "/search-availability", "GET" /*[]postData{},*/, http.StatusOK},
	{"contact", "/search-availability", "GET" /*[]postData{},*/, http.StatusOK},
	{"ms", "/make-reservation", "GET" /*[]postData{},*/, http.StatusOK},

	/* POST */
	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},

	/* this one includes everything (situations) you need to know to write the tests */
	// {"make reservation post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "John"},
	// 	{key: "last_name", value: "Smith"},
	// 	{key: "email", value: "abc@abc.com"},
	// 	{key: "phone", value: "555-55555"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()               // routes is a handler
	ts := httptest.NewTLSServer(routes) // ts is test server
	defer ts.Close()

	for _, e := range theTests {
		// In thsi chapter, we'll focus on GET first

		// if e.method == "GET" {
		// create web client
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
		/*
			} else { // POST
				// a post request as a variable that corresponds to what we're posting in the form
				// an empty variable in the form that is required by the method
				values := url.Values{}
				for _, x := range e.params {
					values.Add(x.key, x.value) // add the key and value to our post
				}
				// call our test clinet, but use PostForm instead of Get
				resp, err := ts.Client().PostForm(ts.URL+e.url, values)
				if err != nil {
					t.Log(err)
					t.Fatal(err)
				}
		*/

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	// Handlers Reservation expects to pull models.Reservation out of the session
	// So we need to have models.Reservation in our TestRepository_Reservation
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	// get our context into the request
	ctx := getCtx(req)
	req = req.WithContext(ctx) // We now have a request that knows about the exta session

	// request record - simulating what we get from the request resposne cycle
	// response cycle: fire a web browser -> hit our website -> get a handler -> pass the request -> get a response writer -> write the response to the web browser
	rr := httptest.NewRecorder() // ResponseRecorder is an implementation of http.ResponseWriter that records its mutations for later inspection in tests.
	// put our reservation in the session with our built context instead of the default context from request
	session.Put(ctx, "reservation", reservation)

	// call the handler resrvation function (but we can't call it directly)
	// turn the handker reservation into a handler function, so we can it directly
	handler := http.HandlerFunc(Repo.Reservation) // now act as a server to run "Repo.Reservation"

	// (m *Repository) Reservation(w http.ResponseWriter, r *http.Request)
	// we call the handler directly with ServeHTTP thae same as we call it on our route in the main funciton
	handler.ServeHTTP(rr, req) // satisfies all the functions to become something that actually acts as a web server
	/*
		We don't even need our routes for this test, so we don't ever call getRoutes at all
		Instead, We built it manually by calling handler.ServeHTTP,
		then passing it my ResponseRecorder and passing it my request,
		which has the necessary sessional information in it.
	*/

	if rr.Code != http.StatusOK {
		t.Errorf("Rservation handler returned wrong response code: %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservations is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-erservation", nil)
	// still need the context with the session header - without this, we can't even test the situation where we can't find the value in the session because there is no session
	ctx = getCtx(req)
	req = req.WithContext(ctx) // now we have a session without reservation because we aren't goona add it
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect { // use StatusTemporaryRedirect because we expects to find that
		t.Errorf("Rservation handler returned wrong response code: %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test with non-existing room
	req, _ = http.NewRequest("GET", "/make-erservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100 // assign the room >2 to simulate the error
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect { // use StatusTemporaryRedirect because we expects to find that
		t.Errorf("Reservation handler returned wrong response code: %d, wanted %d", rr.Code, http.StatusOK)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	// create a POST body
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02") // append the string
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	// Do the same things as TestRepository_Reservation initially, but we can't pass nil but the body to post
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody)) // create a POST to make-reservation
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// Not reqired but a good practice
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // tell the web server that it's a form post

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation) // put it in the handler, and now we can call it

	handler.ServeHTTP(rr, req) // pass something (no need to be ReposeWriter) that satisfies the requirements for being a response writer

	if rr.Code != http.StatusSeeOther { // if everthing is correct, we're expecting StatusSeeOther
		t.Errorf("PostReservation handler returned wrong response code: %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for missing post body (pass the body with nil)
	// copy the same as a bove, change only the body to nil
	req, _ = http.NewRequest("POST", "/make-reservation", nil) // create a POST to make-reservation
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // tell the web server that it's a form post
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation) // put it in the handler, and now we can call it

	handler.ServeHTTP(rr, req) // pass something (no need to be ReposeWriter) that satisfies the requirements for being a response writer

	if rr.Code != http.StatusTemporaryRedirect { // we expect to see StatusTemporaryRedirect
		t.Errorf("PostReservation handler returned wrong response code for missing post body")
	}

	// test for invalid start date
	// copy the same as a bove
	reqBody = "start_date=invalid" // intentionally make the start date invalid
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody)) // create a POST to make-reservation
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") // tell the web server that it's a form post
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation) // put it in the handler, and now we can call it

	handler.ServeHTTP(rr, req) // pass something (no need to be ReposeWriter) that satisfies the requirements for being a response writer

	if rr.Code != http.StatusTemporaryRedirect { // we expect to see StatusTemporaryRedirect according to start_date check from PostReservation
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for invalid end date
	// copy the same as a bove
	reqBody = "start_date=2050-01-01" // intentionally make the start date invalid
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect { // we expect to see StatusTemporaryRedirect according to start_date check from PostReservation
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for invalid room id
	// copy the same as a bove
	reqBody = "start_date=2050-01-01" // intentionally make the start date invalid
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid") // test with something that's not int

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid room id: %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for invalid data
	// copy the same as a bove
	reqBody = "start_date=2050-01-01" // intentionally make the start date invalid
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=J") // test it with the first name less than 3 characters long
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther { // status have to match the error defined in PostReservation
		t.Errorf("PostReservation handler returned wrong response code for invalid data: %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test for failure to insert reservation (InsertReservation) into database
	// copy the same as a bove
	reqBody = "start_date=2050-01-01" // intentionally make the start date invalid
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2") // intentionally generate failure

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect { // status have to match the error defined in PostReservation
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: %d, wanted %d", rr.Code, http.StatusTemporaryRedirect) // not the best comment, but you get the idea
	}

	// test for failure to insert restriction into database
	// copy the same as a bove
	reqBody = "start_date=2050-01-01" // intentionally make the start date invalid
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1000") // intentionally generate failure

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect { // check from handlers.go --> http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		t.Errorf("PostReservation handler failed when trying to fail inserting reservation: %d, wanted %d", rr.Code, http.StatusTemporaryRedirect) // not the best comment, but you get the idea
	}

}

// need to put our reservation variable into the session
// variable --> context --> request
// In order to do that, the context has to have a particular value in there, a key that our session package knows is the session.
func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session")) // ***the header key have to be exactly like this ,so this context knows about the session in order to read to/from and write to the session
	if err != nil {
		log.Println(err)
	}
	return ctx
}
