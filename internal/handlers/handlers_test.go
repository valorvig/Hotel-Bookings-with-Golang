package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
