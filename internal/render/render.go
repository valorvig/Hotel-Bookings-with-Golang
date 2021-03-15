package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/valorvig/bookings/internal/config"
	"github.com/valorvig/bookings/internal/models"
)

// map any functions and parse them to the tempalte
var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

// NewRenderer sets the config fro the template package
// func NewTemplates(a *config.AppConfig) {
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultData allows what data to be available on every page
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// therse will be populated everytime we handle the page
	td.Flash = app.Session.PopString(r.Context(), "flash") // a message (sending to users) appears once and is automatically taken out of the session. PopString puts things in the session until the next time a page is displayed and then it's taken out automatically.
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	// make a decision herer as to whether or not a user is logged in
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1 // logged in
	}

	return td
}

// Template renders templates using html/template
// (Capitalize the name sothat it can be exported)
// It's a ocnvention to not name a function with a given package as the first part of the word - to make it more readable
// func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	// create a variable to hold template cache
	var tc map[string]*template.Template

	// use template cache in development, not in production that needs to rebuild it on every request
	if app.UseCache {
		// get the template cache from the app config
		// read the information from template cache
		tc = app.TemplateCache
	} else { // set UseCache to false in production ro rebuild the template cache
		// rebuild the template cache
		tc, _ = CreateTemplateCache() // get template cache according to each template name with the layout
	}

	// use ok to check whether it exists, not ok if we can't find the template
	t, ok := tc[tmpl] // tempalte from template cache
	// fmt.Println("t: ",t) // &{<nil> 0xc00008a580 0xc0000a6700 0xc0000d0050}
	if !ok {
		// log.Println("can't get template from cache!") // this error is fine (in testing)
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)

	// fill the empty template
	td = AddDefaultData(td, r) // [big] including CSRF token?

	// store the value in buf and don't pass any data (nil)
	err := t.Execute(buf, td) // [Big] We could have wirtten directly to "w" as well instead of having to put the data in "buf" first
	if err != nil {
		log.Fatal(err)
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache creates a template chache as a map
// The purpose is for a structure to hold and look for things quickly
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{} // ready to use parsed templates

	// func Glob(pattern string) (matches []string, err error) - Glob returns the names of all files matching pattern or nil if there is no matching file.
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		fmt.Println("ERROR: pages, err := filepath.Glob(fmt.Sprintf()")
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

/*
// IMPORTANT!
// use this directory or else all the paths will be mismatch
go run cmd/web/*.go
*/
