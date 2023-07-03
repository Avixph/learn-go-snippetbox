package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/google/uuid"
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
	// convert it to uuid using the uuid.Parse() func. If it's value is less than nil, we return a 404 page not found response.
	id, err := uuid.Parse(r.URL.Query().Get("id"))
	if err != nil {
		// Use the notFound() helper
		app.notFound(w)
		return
	}

	// Use the SnippetModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found, then
	// return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the fmt.Fprintf() func to write the snippet data as a plain-text
	// HTTP response body.
	fmt.Fprintf(w, "%+v", snippet)
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

	// Create a few vars holding dummy data. We'll remove these later on
	// during the build.
	title := "0 snails"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expireVal := 7

	// Pass the data to the SnippetModel.Insert() method, receive the ID of
	// the new record back.
	id, err := app.snippets.Insert(title, content, expireVal)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// w.Write([]byte("Create a new snippet..."))
	// Rediect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%v", id), http.StatusSeeOther)
}
