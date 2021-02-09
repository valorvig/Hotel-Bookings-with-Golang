package handlers

import (
	"net/http"

	"github.com/valorvig/bookings/pkg/config"
	"github.com/valorvig/bookings/pkg/models"
	"github.com/valorvig/bookings/pkg/render"
)

// we can't have the struct model templatedata here since it's gpnna create import cycle error

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type (Repository pattern)
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
// With this repository receiver, now this handler can access everything inside the repository, especially the AppConfig
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// grab the IP address of the person visiting the site and store it in the home page session
	remoteIP := r.RemoteAddr // IPv4 or IPv6 address
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again."

	// hold user's IP address in the session
	// the value is empty if there is nothing in the session named "remote_ip"
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")

	// after accessing the "session" from config, can do anything from it
	// m.App.Session.

	stringMap["remote_ip"] = remoteIP // Ex. [::1]:18107 - "::1" is the loopback address in ipv6, equal to 127.0.0.1 in ipv4.

	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}
