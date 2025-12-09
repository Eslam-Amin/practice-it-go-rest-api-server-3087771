package main

import (
	"net/http"

	"example.com/backend"
)

// Main is the entry point of the program. It sets up the backend application
// and starts the HTTP server.
func main() {
	app := &backend.App{
		Port: ":8080",
	}

	app.Initialize()
	http.Handle("/", app.Router)
	app.Run()
}
