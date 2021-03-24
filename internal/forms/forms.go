package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct, embeds a url.Values object
// type embedding based on url.Values
type Form struct {
	url.Values
	Errors errors
}

// Valid turns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New initializes a form struct
func New(data url.Values) *Form {
	// [Big] kind of creating a wrapper of url.Values with adding additional "errors"
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
// use this to avoid the hard codes for many "form.Has("first_name", r)"
func (f *Form) Required(fields ...string) { // ex. first_name, last_name, email
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" { // "" if no value
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has checks if the form field is in post and not empty.
// Has also checks whether ot not submitted form data includes a certain field
// ex. checkbox and text are handled diffrently.
// func (f *Form) Has(field string, r *http.Request) bool {
func (f *Form) Has(field string) bool {
	// x := r.Form.Get(field) // we should check the one associates with our receiver f not the request r
	x := f.Get(field)
	if x == "" {
		// f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	// x := r.Form.Get(field)
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
