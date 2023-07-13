package main

import (
	"errors"
	"fmt"
	"strconv"

	"net/http"

	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

// Define a home handler func that writes a byte slice containing
// "Hello from Snippetbox!" as the response body.
// Change the signature if the home handler so it is defined as a method
// against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because httprouter matches the "/" path exactly, we don't need the manual
	// check of `if r.URL.Path != "/"` from the handler.

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
	// When httprouter is parsing a request, the values of any name parameters
	// will be stored in the request context. We can use the ParamsFromContext()
	// func to retrieve a slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(r.Context())

	// We can use the ByName() method to get the value of the "id" named parameter
	// from the slice and validate it as normal.
	id, err := uuid.Parse(params.ByName("id"))
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

// Define snippetCreateForm handler func, which for now returns a placeholder.
func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	TemplData := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.html", TemplData)
}

// Define snippetCreate handler func
// Change the signature if the snippetCreate handler so it is defined as a
// method against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Call the r.ParseForm() method to add any data in POST request bodies to the
	// r.PostForm map (also works in the same way for PUT and PATCH requests). If
	// there are any errors, we use our app.ClientError() helper to send a 400 Bad
	// Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// Use the r.PostForm.Get() method to retrieve the title and content from
	// the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The r.PostForm.Get() method always returns the form dat as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an iteger. So we need to manually covert
	// the form data to an integer using strcov.Atoi(), and we send a 400 BAD
	// REQUEST response if the conversion fails.
	expireVal, err := strconv.Atoi(r.PostForm.Get("expires value"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	// Checking if the request method is a POST is now superfluous, because
	// this is done by httprouter automatically.

	// Create a few vars holding dummy data. We'll remove these later on
	// during the build.
	// title := "0 snails"
	// content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	// expireVal := 7

	// Pass the data to the SnippetModel.Insert() method, receive the ID of
	// the new record back.
	id, err := app.snippets.Insert(title, content, expireVal)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// w.Write([]byte("Create a new snippet..."))
	// Rediect the user to the relevant page for the snippet.
	// Update the redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", id), http.StatusSeeOther)
}
