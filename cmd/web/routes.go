package main

import "net/http"

// The routes() method returns a servermux containing our application routes.
func (app *application) routes() *http.ServeMux {
	// Use the http.NewServerMux() func to initialize a new servermux.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static"
	// directory. Note that the path given to the http.Dir func is relative
	// to the project directory root.
	// Use the mux.handle() func to register the file server as the handler
	// for all URL paths that start with "/static/". For matching paths, we
	// strip the "/static/" prefix before the request reaches the file
	// server.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the home, snippetView and snippetCreate funcs as handlers
	// for the corrisponding URL patrerns with the servermux.
	// Swap the route declearations to use the application struct's methods
	// as the handler func.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
