package main

import (
	"errors"
	"fmt"
	"strconv"

	"net/http"

	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/Avixph/learn-go-snippetbox/internal/validator"
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
	templData := app.newTemplateData(r)
	templData.Snippets = snippets

	// Use the new render helper.
	app.render(w, http.StatusOK, "home.html", templData)
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
	templData := app.newTemplateData(r)
	templData.Snippet = snippet

	// Use the new render helper.
	app.render(w, http.StatusOK, "view.html", templData)
}

// Define snippetCreateForm handler func, which for now returns a placeholder.
func (app *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	templData := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	templData.Form = snippetForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.html", templData)
}

// Define a snippetForm struct to represent the form data and validation errors
// for the form fields. Note: all the struct fields are deliberately exported
// (i.e. start with a capital letter). This is because struct fields must be
// exported in order to be read by the html/template package when rendering
// the template.
// Embed the Validator type which will allow the snippetForm to "inherit"
// all the fields and methods of our validator type (including the
// FieldErrors field).
type snippetForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors validator.Validator
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
		return
	}

	// The r.PostForm.Get() method always returns the form dat as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an iteger. So we need to manually covert
	// the form data to an integer using strcov.Atoi(), and we send a 400 BAD
	// REQUEST response if the conversion fails.
	expireVal, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create an instance of the snippetForm struct containing the values from
	// the form and an empty map for any validation errors.
	form := snippetForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expireVal,
	}

	// Since the validator type is embedded by the snippetForm struct, we can call CheckField() directly on iy to execute our validation checks. CheckField() will add the provided key and error message to the FieldErrors map if the check does not evaluate to true.
	form.FieldErrors.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank!")
	form.FieldErrors.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long!")
	form.FieldErrors.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank!")
	form.FieldErrors.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365!")

	// Use the Valid() method to see if any of the check failed. If they did,
	// then re-render the template passing in the form in the same way as before.
	if !form.FieldErrors.Valid() {
		templData := app.newTemplateData(r)
		templData.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", templData)
		return
	}

	// Pass the data to the SnippetModel.Insert() method, receive the ID of
	// the new record back.
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// w.Write([]byte("Create a new snippet..."))
	// Rediect the user to the relevant page for the snippet.
	// Update the redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", id), http.StatusSeeOther)
}
