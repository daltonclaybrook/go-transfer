package main

import (
	"github.com/daltonclaybrook/swerve"
	"github.com/daltonclaybrook/swerve/middle"
)

func main() {
	server := swerve.NewServer()
	server.AddMiddleware(middle.CORS{Origin: "*", Methods: "POST, GET, OPTIONS", Headers: "*"})
	server.AddControl(NewTransfer())
	server.Start()
}
