package models

import "github.com/valorvig/bookings/internal/forms"

// TemplateData holds data sent from handlers to templates
// We may want ot sent any data that we can't decide yet to RenderTemplate, so we create a struct to hold such data
// it will only exists to be imported by packages other than itself
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string // flash message
	Warning   string
	Error     string
	Form      *forms.Form // whether a page has a form or not, the form object is available here
}
