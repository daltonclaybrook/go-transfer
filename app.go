package main

import (
	"github.com/daltonclaybrook/transfer/controller"
	"github.com/daltonclaybrook/transfer/server"
)

func main() {
	server := server.WebServer{}
	server.RegisterController(controller.NewTransfer())
	server.Start()
}
