package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/justinas/nosurf"
)

func secureHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note this is split across multiple lines for readability.
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred func that will always be run in the event of a panic as Go
		// unwinds the stack.
		defer func() {
			// Use the builtin recover func to check if there has been a panic or not.
			// If there is...
			if err := recover(); err != nil {
				// Set a "Connection: close" header to the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500 Internal Server
				// Response.
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)

	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the user is not authenticated, redirect them to the login page and
		// return from the middleware chain so that no subsequent handlers in the
		// chain are executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		}

		// Else, set the "Cache-Control: no-store" header so that pages required
		// authentication are not stored in the user's browser cache (or other
		// intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		// Finally call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// Create a NoSurf middleware func which uses a customized CSRF coockie with the
// Secure, Path and HttpOnly attributes set.
func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the authenticatedUserID value from the session using the
		// GetString() method. This will the string value if no
		// "authenticatedUserID" value is in the session -- in which case we call
		// the next handler in the chain as normal and return.
		id := app.sessionManager.GetString(r.Context(), "authenticatedUserID")
		if id == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Else, check if a user with that id exists in our database.
		exists, err := app.users.Exists(uuid.MustParse(id))
		if err != nil {
			app.serverError(w, err)
			return
		}

		// If a matching user is found, we know that the request is coming from an
		// authenticated user who exists in our database. We create a new copy of
		// the request (with an isAuthenticatedContextKey value of true in the request
		// context) and assign it to r.
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedCOntextKey, true)
			r = r.WithContext(ctx)
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})
}
