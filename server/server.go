package server

import (
	"fmt"
	"github.com/daltonclaybrook/go-transfer/controller"
	"github.com/daltonclaybrook/go-transfer/middle"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

// WebServer is used to create and start a server.
type WebServer struct {
	server     *http.Server
	router     *mux.Router
	middleware []middle.Middle
}

// NewWebServer returns a new initialized instance of WebServer.
func NewWebServer() *WebServer {
	ws := &WebServer{}
	ws.middleware = make([]middle.Middle, 0)
	ws.router = mux.NewRouter()
	ws.router.HandleFunc("/", sendUnhandled)
	http.Handle("/", ws.router)
	return ws
}

// RegisterController registers a request handler with the WebServer.
func (ws *WebServer) RegisterController(c controller.Controller) {
	routes := c.Routes()
	for _, route := range routes {
		for _, handler := range route.Handlers {
			fmt.Printf("path: %v, method: %v\n", route.Path, handler.Method)

			// temp hack
			newHandler := ws.middleware[0].Handle(handler.Handler)
			ws.router.HandleFunc(route.Path, newHandler).Methods(handler.Method)
		}
	}
}

func (ws *WebServer) RegisterMiddleware(m middle.Middle) {
	ws.middleware = append(ws.middleware, m)
}

// Start starts the server.
func (ws *WebServer) Start() {
	ws.setupServer()
	ws.server.ListenAndServe()
}

/*
Private
*/

func (ws *WebServer) setupServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ws.server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	ws.server.ErrorLog = log.New(os.Stdout, "err: ", 0)
}

func sendUnhandled(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "Method \"%v\" is not supported by this route.", r.Method)
}
