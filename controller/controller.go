package controller

import (
	"net/http"
)

// Route describes an endpoint.
type Route struct {
	Path     string
	Handlers []Handler
}

// Handler describes functions mapped to http methods.
type Handler struct {
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}

// Controller handles routes.
type Controller interface {
	Routes() []Route
}
