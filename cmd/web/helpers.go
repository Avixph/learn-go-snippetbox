package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
)

// The serverError helper writes an error message and stack trace to the
// errorLog, then sends a generic 500 Internal Server Error response to the
// user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	ise := http.StatusInternalServerError
	http.Error(w, http.StatusText(ise), ise)
}

// The clientError helper sends a specific status code and corresponding
// description to the user. We'll use this later to send responses like 400
// "Bad Request" when there's a problem with the request the user sent.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// We'll also implement a notFound helper. This is simply a convenience
// wrapper around clientError which sends a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	nf := http.StatusNotFound
	app.clientError(w, nf)
}

// The getEnvVariables() helper reads the .env file and returns the requested
// key value
func getEnvVariables(key string) string {
	// Use the godotenv.Load() method to load the env files.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("No .env file found!")
	}
	return os.Getenv(key)
}
