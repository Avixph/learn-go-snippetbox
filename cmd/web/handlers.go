package main

import (
	"errors"
	"fmt"

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

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containg
	// the 'default' data (which for now is just the current year), and add
	// the snippet slice to it.
	TemplData := app.newTemplateData(r)
	TemplData.Snippets = snippets

	// Use the new render helper.
	app.render(w, http.StatusOK, "home.html", TemplData)
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

	// Call the newTemplateData() helper to get a templateData struct containg
	// the 'default' data (which for now is just the current year), and add
	// the snippet slice to it.
	TemplData := app.newTemplateData(r)
	TemplData.Snippet = snippet

	// Use the new render helper.
	app.render(w, http.StatusOK, "view.html", TemplData)
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
