package controller

import (
	// "encoding/json"
	"fmt"
	// "github.com/daltonclaybrook/web-app/model"
	"github.com/gorilla/mux"
	"net/http"
)

// User is the controller for all /user endpoints.
type User struct{}

// Routes describes all endpoints handled by the User Controller.
func (uc *User) Routes() []Route {
	return AllRoutesFromHandler("user", uc)
}

/*
Handlers
*/

func (uc *User) create(w http.ResponseWriter, r *http.Request) {
	// var user *model.User
	// json.Unmarshal()

	fmt.Fprintln(w, "create")
}

func (uc *User) find(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "find")
}

func (uc *User) findOne(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Fprintf(w, "findOne with id: %v\n", id)
}

func (uc *User) update(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "update")
}

func (uc *User) delete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "delete")
}
