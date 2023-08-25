package main

import (
	"net/http"

	"github.com/Avixph/learn-go-snippetbox/ui"
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

	// Take the ui.Files embeded filesytem and convert it to a a http.FS type so
	// that it satisfies the http.FileSystem interface. We then pass that to the
	// http.FileServer() func to create the file server handler.
	// Use the router.Handler() func to register the file server as the handler
	// for all URL paths that start with "/static/".
	// Our static files are contained in the "static" folder of the embeded ui.
	// Files (ex: The CSS stylesheet is located at "/static/css/main.css" which
	// means that we no longer need to strip the prefix from the request URL. Any \
	// request that starts with "/static/" can just be passed directly to the
	// fileserver and corresponding static file will be served.)
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Add a GET /ping route.
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Create a middleware chain containing the middleware specific to our
	// unprotected application routes using the "dynamic" middleware chain.
	// Use the noSurf and authenticate middleware on all our 'dynamic' routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Register the home, snippetView and snippetCreate funcs as handlers for the
	// corrisponding URL patrerns with the serverrouter. Swap the route
	// declearations to use the application struct's methods as the handler func.
	// Update routes to use the new dynamic middleware chain followed by the
	// appropriate handler func. Note: Because alice ThenFunc() method returns
	// an http.Handler() instead of an http.HandlerFunc() we also need to switch
	// to registering the route using the router.Hanler() method.
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	// Add the About route.
	router.Handler(http.MethodGet, "/about", dynamic.ThenFunc(app.about))
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
	router.Handler(http.MethodGet, "/account/view", protected.ThenFunc(app.accountView))
	router.Handler(http.MethodGet, "/account/password/update", protected.ThenFunc(app.userPasswordUpdateForm))
	router.Handler(http.MethodPost, "/account/password/update", protected.ThenFunc(app.userPasswordUpdate))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogout))

	// Create a middleware chain containing our 'standard' middleware (app.recoverPanic,
	// app.logRequest, secureHeader) which will be used for every request received.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeader)

	// Return the 'standard' middleware chain followed by the httprouter
	return standard.Then(router)
}
