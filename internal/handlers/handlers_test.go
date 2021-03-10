package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	name               string     // name of the test
	url                string     // the path which is matched by our roots
	method             string     // GET or POST
	params             []postData // things that are being posted and that's a form can have more than 1 input
	expectedStatusCode int        // what kind of response code we're getting back from the web server (200, 404, etc.)
}{
	/*
		{"home", "/", "GET", []postData{}, http.StatusOK},
		{"about", "/about", "GET", []postData{}, http.StatusOK},
		{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
		{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
		{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
		{"contact", "/search-availability", "GET", []postData{}, http.StatusOK},
		{"ms", "/make-reservation", "GET", []postData{}, http.StatusOK},
		{"post-search-avail", "/search-availability", "POST", []postData{
			{key: "start", value: "2020-01-01"},
			{key: "end", value: "2020-01-02"},
		}, http.StatusOK},
		{"post-search-avail-json", "/search-availability-json", "POST", []postData{
			{key: "start", value: "2020-01-01"},
			{key: "end", value: "2020-01-02"},
		}, http.StatusOK},
		{"make reservation post", "/make-reservation", "POST", []postData{
			{key: "first_name", value: "John"},
			{key: "last_name", value: "Smith"},
			{key: "email", value: "abc@abc.com"},
			{key: "phone", value: "555-55555"},
		}, http.StatusOK},
	*/
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()               // routes is a handler
	ts := httptest.NewTLSServer(routes) // ts is test server
	defer ts.Close()

	for _, e := range theTests {
		// In thsi chapter, we'll focus on GET first
		if e.method == "GET" {
			// create web client
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
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

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
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
	rr := httptest.NewRecorder()
	// put our reservation in the session with our built context instead of the default context from request
	session.Put(ctx, "reservation", reservation)

	// call the handler resrvation function (but we can't call it directly)
	// turn the handker reservation into a handler function, so we can it directly
	handler := http.HandlerFunc(Repo.Reservation) // now act as a server

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
