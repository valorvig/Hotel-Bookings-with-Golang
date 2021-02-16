package render

import (
	"bytes"
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

// NewTemplates sets the config fro the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// AddDefaultData allows what data to be available on every page
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
// (Capitalize the name sothat it can be exported)
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// create a variable to hold template cache
	var tc map[string]*template.Template

	// use template cache in development, not in production that needs to rebuild it on every request
	if app.UseCache {
		// get the template cache from the app config
		// read the information from template cache
		tc = app.TemplateCache
	} else { // set UseCache to false in production ro rebuild the template cache
		// rebuild the template cache
		tc, _ = CreateTemplateCache()
	}

	// use ok to check whether it exists, not ok if we can't find the template
	t, ok := tc[tmpl] // tempalte from template cache
	// fmt.Println("t: ",t) // &{<nil> 0xc00008a580 0xc0000a6700 0xc0000d0050}
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	// store the value in buf and don't pass any data (nil)
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template chache as a map
// The purpose is for a structure to hold and look for things quickly
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{} // ready to use parsed templates

	// func Glob(pattern string) (matches []string, err error) - Glob returns the names of all files matching pattern or nil if there is no matching file.
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		// func Base(path string) string - Base returns the last element of path. Trailing path separators are removed before extracting the last element.
		name := filepath.Base(page) // name is about.page.tmpl or home.page.tmpl

		// fmt.Println("Page is currently", page) // Ex. Page is currently templates\about.page.tmpl

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		// fmt.Println("ts1: ", ts)
		// Home --> &{<nil> 0xc0000205c0 0xc00013c200 0xc00004e1e0}
		// About --> &{<nil> 0xc00008a180 0xc0000a6100 0xc0000ce000}
		if err != nil {
			// fmt.Println("ERROR: ts, err := template.New(name).Funcs(functions).ParseFiles(page)")
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		// fmt.Println("matches: ", matches) // [templates\base.layout.tmpl]
		if err != nil {
			// fmt.Println("ERROR:	matches, err := filepath.Glob(\"./templates/*.layout.tmpl\")")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			// fmt.Println("ts2: ", ts)
			// Home --> &{<nil> 0xc0000205c0 0xc00013c200 0xc00004e1e0}
			// About --> &{<nil> 0xc00008a180 0xc0000a6100 0xc0000ce000}
			if err != nil {
				// fmt.Println("ERROR: ts, err = ts.ParseGlob(\"./templates/*.layout.tmpl\")")
				return myCache, err
			}
		}

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
