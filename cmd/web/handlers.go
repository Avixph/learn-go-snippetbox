package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Define a home handler func that writes a byte slice containing
// "Hello from Snippetbox!" as the response body.
// Change the signature if the home handler so it is defined as a method
// against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If not, use
	// the app.notFound() func to send a 404 respond to the client.
	// Importantly, we then return from the handler. If we don't return, the
	// handler would keep executing and also write "Hello from Snippetbox!" message
	rup := r.URL.Path
	if rup != "/" {
		// Use the notFound() helper
		app.notFound(w)
		return
	}

	// Initialize a slice containing the paths to the two templates. It's
	// important to note that the file containing our base template must be
	// the "first" file in the slice.
	templateFiles := []string{
		"./ui/html/base.html",
		"./ui/html/components/nav.html",
		"./ui/html/pages/home.html",
		"./ui/html/components/footer.html",
	}

	// Use the template.ParseFiles() func to read the template files and
	// store the templates in a template set. If there's an err, we log a
	// detailed err message and use the http.Error() func to send a generic
	// 500 Interanl Server Err response to the user.
	ts, err := template.ParseFiles(templateFiles...)
	if err != nil {
		// Because the home handler func is now a method against application
		// it can access it's feilds, including the error loger. We'll write
		// the log message to this instead of the standard logger.
		// Use ther serverError() helper
		app.serverError(w, err)
		return
	}

	// We then use the ExecuteTemplate() method to write the content of the
	// "base" template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		// Also update the code hre to use the error logger from the
		// application struct.
		// Use the app.serverError() helper
		app.serverError(w, err)
	}
}

// Define a snippetView handler func
// Change the signature if the snippetView handler so it is defined as a
// method against *application.
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() func. If it can't be
	// converted to an integer, or it's value is less than 0, we return a 404
	// page not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		// Use the notFound() helper
		app.notFound(w)
		return
	}

	// Use the fmt.Fprintf() func to interpolate the id value with our
	// response and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Displaying a specific snippet with ID# %d...", id)
}

// Define snippetCreate handler func
// Change the signature if the snippetCreate handler so it is defined as a
// method against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	rm := r.Method
	if rm != http.MethodPost {
		// Use the Header().Set() method to add a 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second rarameter is the header value.
		w.Header().Set("Allow", http.MethodPost)

		// If it's not, use the clientError() helper to send a 405 status code
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
