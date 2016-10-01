package main

import (
	"github.com/daltonclaybrook/go-transfer/controller"
	"github.com/daltonclaybrook/go-transfer/server"
)

func main() {
	server := server.WebServer{}
	server.RegisterController(controller.NewTransfer())
	server.Start()
}
