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
	// Update the pattern for the route for the static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// Create a middleware chain containing the middleware specific to our
	// unprotected application routes using the "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// Register the home, snippetView and snippetCreate funcs as handlers for the
	// corrisponding URL patrerns with the serverrouter. Swap the route
	// declearations to use the application struct's methods as the handler func.
	// Update routes to use the new dynamic middleware chain followed by the
	// appropriate handler func. Note: Because alice ThenFunc() method returns
	// an http.Handler() instead of an http.HandlerFunc() we also need to switch
	// to registering the route using the router.Hanler() method.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignupForm))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLoginForm))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLogin))

	// Create a protected (authenticated) middleware chain containing the
	// middleware specific to our "protected" middleware chain which includes the
	// requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreateForm))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogout))

	// Create a middleware chain containing our 'standard' middleware (app.recoverPanic,
	// app.logRequest, secureHeader) which will be used for every request received.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	// Return the 'standard' middleware chain followed by the httprouter
	return standard.Then(router)
}
