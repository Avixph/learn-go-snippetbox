package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The routes() method returns a http.Handler containing our application routes.
func (app *application) routes() http.Handler {
	// Initialize a new httprouter router
	router := httprouter.New()

	// Create a handler func that wraps our notFound() helper, and then assign it
	// as the custom handler for 404 Not Found responses. We can also set a
	// custom handler for 405 Method Not Allowed responses by setting router.
	// MethodNotAllowed in the same way too.
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// Create a file server which serves files out of the "./ui/static"
	// directory. Note that the path given to the http.Dir func is relative
	// to the project directory root.
	// Use the router.Handler() func to register the file server as the handler
	// for all URL paths that start with "/static/". For matching paths, we
	// strip the "/static/" prefix before the request reaches the file
	// server.
	// Update the pattern for the route for the static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Register the home, snippetView and snippetCreate funcs as handlers for the
	// corrisponding URL patrerns with the serverrouter. Swap the route
	// declearations to use the application struct's methods as the handler func.
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreateForm)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreate)

	// Create a middleware chain containing our 'standard' middleware (app.recoverPanic,
	// app.logRequest, secureHeader) which will be used for every request received.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	// Return the 'standard' middleware chain followed by the httprouter
	return standard.Then(router)
}
