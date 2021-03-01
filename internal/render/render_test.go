package render

import (
	"net/http"
	"testing"

	"github.com/valorvig/bookings/internal/models"
)

// the applicaion funciton that we want to test first is AddDefaultData.
// AddDefaultData takes "td *models.TemplateData" and "r *http.Request".

// we want to test AddDeafaultData first - requires td and r as inputs
func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// let's create something to test in the session
	session.Put(r.Context(), "flash", "123")

	// let's try calling AddDefaultData after we've provided both td and r.
	// With getSession above, when we call AddDefaultData, we can call it with an empty td (a pointer to template data) and the request r (that now has session information)
	result := AddDefaultData(&td, r)

	// without the session data, it would definitely fail
	// We've put data into "flash", we can now specifically check "result.Flash" instead of "result" alone.
	if result.Flash != "123" {
		t.Error("flash value of 123 not found in sesison")
	}
}

func TestTemplate(t *testing.T) {
	// the path for the test is two levels
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	// r is a poniter to http.Reequest
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// we've create the necessary variable to satisfy the requirements for a ResponseWriter
	var ww myWriter

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{
		// currently, we don't have a Writer and a Response
	})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Template(&ww, r, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

// not a test fucntion but our own define within this package
func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil) // making request to url trying to get a request
	// NewRequest is a built-in func and it's not going to give error. However, checking error is a good habit.
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session")) // the second parameter makes it an active session
	r = r.WithContext(ctx)                                // put the context back into the request

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplatecache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

/*
cd internal/render

varut@BIG MINGW64 /d/big_golang/trevor/bookings/internal/render (main)
$ go test -v

varut@BIG MINGW64 /d/big_golang/trevor/bookings/internal/render (main)
$ go test -cover
2021/02/21 20:16:49 can't get template from cache!
PASS
coverage: 81.0% of statements
ok      github.com/valorvig/bookings/internal/render    0.354s

// linux command to provide the test coverage
// create a file called coverage.out
// On Mac and Linux, && checks if the prior command succeeds, then run the next command
// next command needs to cover -html equals the name of the file
$ go test -coverprofile=coverage.out && go tool cover -html=coverage.out

*/
