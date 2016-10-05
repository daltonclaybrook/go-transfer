package controller

import (
	"github.com/daltonclaybrook/go-transfer/middle"
)

// Route describes an endpoint.
type Route struct {
	Path     string
	Handlers []Handler
}

// Handler describes functions mapped to http methods.
type Handler struct {
	Method  string
	Handler middle.ContextFunc
}

// Controller handles routes.
type Controller interface {
	Routes() []Route
}
