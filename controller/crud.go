package controller

import (
	"fmt"
	"net/http"
)

const (
	// Create a model.
	Create = iota
	// Find models.
	Find = iota
	// FindOne model.
	FindOne = iota
	// Update a model.
	Update = iota
	// Delete a model.
	Delete = iota
)

// CRUDRoute is used to express a typical CRUD operation.
type CRUDRoute struct {
	Op      int
	Handler func(w http.ResponseWriter, r *http.Request)
}

// CRUDHandler defines a type which handles all CRUD methods.
type CRUDHandler interface {
	create(w http.ResponseWriter, r *http.Request)
	find(w http.ResponseWriter, r *http.Request)
	findOne(w http.ResponseWriter, r *http.Request)
	update(w http.ResponseWriter, r *http.Request)
	delete(w http.ResponseWriter, r *http.Request)
}

// AllRoutesFromHandler is a convenience method for subscribing to all routes.
func AllRoutesFromHandler(model string, handler CRUDHandler) []Route {
	crud := []CRUDRoute{
		CRUDRoute{Create, handler.create},
		CRUDRoute{Find, handler.find},
		CRUDRoute{FindOne, handler.findOne},
		CRUDRoute{Update, handler.update},
		CRUDRoute{Delete, handler.delete},
	}
	return RoutesFromCRUD(model, crud)
}

// RoutesFromCRUD transforms CRUDROutes to Routes expected by the server.
func RoutesFromCRUD(model string, crud []CRUDRoute) []Route {

	m := make(map[string]*Route)
	for _, r := range crud {

		getRoute := func(pattern string) *Route {
			route := m[pattern]
			if route == nil {
				route = &Route{Path: pattern}
				m[pattern] = route
			}
			return route
		}

		switch r.Op {
		case Create:
			route := getRoute(fmt.Sprintf("/%v", model))
			route.Handlers = append(route.Handlers, Handler{"post", r.Handler})
		case Find:
			route := getRoute(fmt.Sprintf("/%v", model))
			route.Handlers = append(route.Handlers, Handler{"get", r.Handler})
		case FindOne:
			route := getRoute(fmt.Sprintf("/%v/{id:[0-9]+}", model))
			route.Handlers = append(route.Handlers, Handler{"get", r.Handler})
		case Update:
			route := getRoute(fmt.Sprintf("/%v/{id:[0-9]+}", model))
			route.Handlers = append(route.Handlers, Handler{"patch", r.Handler})
		case Delete:
			route := getRoute(fmt.Sprintf("/%v/{id:[0-9]+}", model))
			route.Handlers = append(route.Handlers, Handler{"delete", r.Handler})
		}
	}

	retRoutes := make([]Route, 0, len(m))
	for _, value := range m {
		retRoutes = append(retRoutes, *value)
	}

	return retRoutes
}
