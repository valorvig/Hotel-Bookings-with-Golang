package models

import "github.com/valorvig/bookings/internal/forms"

// TemplateData holds data sent from handlers to templates
// We may want ot sent any data that we can't decide yet to RenderTemplate, so we create a struct to hold such data
// it will only exists to be imported by packages other than itself
// This struct holds information that is available to every single page of our website.
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string      // flash message to notie
	Warning         string      // to notie
	Error           string      // to notie
	Form            *forms.Form // whether a page has a form or not, the form object is available here
	IsAuthenticated int         // >0 = logged in, 0 = not logged in
}
