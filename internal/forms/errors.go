package forms

type errors map[string][]string // we might have more than one error for a given field in a form

// Add adds error message for a given form field
func (e errors) Add(field, message string) {
	// field refers to name="first_name", "last_name", and "email"
	e[field] = append(e[field], message)
}

// Get returns the first error message
func (e errors) Get(field string) string {
	// a map of whatever we find in the index of field
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	// return the first index on the error string - to see if the given field has an error
	return es[0] // display the first error depending on when Get first executes
}
