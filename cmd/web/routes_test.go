package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/valorvig/bookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux: // routes returns mux, and mux is *chi.Mux (when you hover your cursor at it)
		// do nothing; test passed
	default:
		t.Error(fmt.Sprintf("type is not *chi.Mux, type is %T", v))
	}
}
