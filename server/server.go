package server

import (
	"fmt"
	"github.com/daltonclaybrook/transfer/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// WebServer is used to create and start a server.
type WebServer struct {
	server *http.Server
	router *mux.Router
}

// RegisterController registers a request handler with the WebServer.
func (ws *WebServer) RegisterController(c controller.Controller) {
	if ws.router == nil {
		ws.router = mux.NewRouter()
		ws.router.HandleFunc("/", sendUnhandled)
		http.Handle("/", ws.router)
	}

	routes := c.Routes()
	for _, route := range routes {
		for _, handler := range route.Handlers {
			fmt.Printf("path: %v, method: %v\n", route.Path, handler.Method)
			ws.router.HandleFunc(route.Path, handler.Handler).Methods(handler.Method)
		}
	}
}

// Start starts the server.
func (ws *WebServer) Start() {
	ws.setupServer()
	fmt.Println("listening on localhost:8080")
	ws.server.ListenAndServe()
}

/*
Private
*/

func (ws *WebServer) setupServer() {
	ws.server = &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	ws.server.ErrorLog = log.New(os.Stdout, "err: ", 0)
}

func sendUnhandled(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Method \"%v\" is not supported by this route.", r.Method)
}
