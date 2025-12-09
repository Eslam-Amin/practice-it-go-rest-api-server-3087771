package main

import (
	"net/http"

	"example.com/backend"
)

func main() {
	app := &backend.App{
		Port: ":8080",
	}

	app.Initialize()
	http.Handle("/", app.Router)
	app.Run()
}
