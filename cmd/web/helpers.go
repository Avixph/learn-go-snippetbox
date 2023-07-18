package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-playground/form/v4"
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

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropiate set from the cache based on the page name (ex:
	// 'home.html'). If no entrry exists in the cache withthe provided name,
	// then create a new err and call the serverError() helper method.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer.
	buff := new(bytes.Buffer)

	// Write the template to the buffer, instead of a straight to the http.
	// ResponseWriter. If an err occurs, call the serverError() helper and
	// then return.
	err := ts.ExecuteTemplate(buff, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without err, then we are safe
	// to go ahead and write out the provided HTTP status code ('200 OK', '400
	// Bad Request' etc) to the http.ResponseWriter.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function
	// that takes an io.Writer.
	buff.WriteTo(w)
}

// Create a newTemplateData() helper, which returns a pointer to a
// templateData struct initialized with the current year. Note: we're not
// using the *http.Request parameter here at the moment, but we will later.
// add the flash message to the template data, if one exists.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		Flash:       app.sessionManager.PopString(r.Context(), "flash"),
	}
}

// Create a decodePostForm() helper method. THe second parameter here, dst, is
// the target destination that we want to decode the form data into.
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// createSnippet handler.
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try yo use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError. We use
		// errors.As() to check for this and raise a panic rather than returning
		// the err.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errs, we return tham as normal
		return err
	}
	return nil
}
