package main

import (
	"github.com/daltonclaybrook/go-transfer/controller"
	"github.com/daltonclaybrook/go-transfer/middle"
	"github.com/daltonclaybrook/go-transfer/server"
)

func main() {
	server := server.NewWebServer()
	server.RegisterMiddleware(middle.CORS{Origin: "*", Methods: "POST, GET, OPTIONS"})
	server.RegisterController(controller.NewTransfer())
	server.Start()
}
