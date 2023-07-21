package main

import (
	"errors"
	"fmt"

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

	// Use the PopString() method to retrieve the value for the "flash" key. The
	// method also deletes the key and value from the session data, so that it
	// acts like a one-time fetch. If there is no matching key in the session
	// data this will return the empty string.
	// flash := app.sessionManager.PopString(r.Context(), "flash")

	// Call the newTemplateData() helper to get a templateData struct containg
	// the 'default' data (which for now is just the current year), and add
	// the snippet slice to it.
	templData := app.newTemplateData(r)
	templData.Snippet = snippet

	// Pass the flash message to the template.
	// templData.Flash = flash

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
// Add struct tags that tell the decoder how to map HTML form values into
// different struct fields. (i.e. Bellow we're telling the decoder to store
// the value from the HTML form inputs with the name "title" in the Title
// field. The struct tag `form:"-"` tells the decoder to completely ignore a
// field during decoding.)
type snippetForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

// Define snippetCreate handler func
// Change the signature if the snippetCreate handler so it is defined as a
// method against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Declare a new empty instance of the snippetForm struct.
	var form snippetForm

	// Call the decodePostForm() helper, passing in the current request and
	// *a pointer* to our snippetForm struct. This will essentially fill our
	// struct with the relevant values from the HTML form.If there is a problem,
	// we return a 400 Bed Request response to the client.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// The r.PostForm.Get() method always returns the form dat as a *string*.
	// However, we're expecting our expires value to be a number, and want to
	// represent it in our Go code as an iteger. So we need to manually covert
	// the form data to an integer using strcov.Atoi(), and we send a 400 BAD
	// REQUEST response if the conversion fails.
	// expireVal, err := strconv.Atoi(r.PostForm.Get("expires"))
	// if err != nil {
	// 	app.clientError(w, http.StatusBadRequest)
	// 	return
	// }

	// Create an instance of the snippetForm struct containing the values from
	// the form and an empty map for any validation errors.
	// form := snippetForm{
	// 	Title:   r.PostForm.Get("title"),
	// 	Content: r.PostForm.Get("content"),
	// 	Expires: expireVal,
	// }

	// Since the validator type is embedded by the snippetForm struct, we can call CheckField() directly on iy to execute our validation checks. CheckField() will add the provided key and error message to the FieldErrors map if the check does not evaluate to true.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank!")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long!")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank!")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365!")

	// Use the Valid() method to see if any of the check failed. If they did,
	// then re-render the template passing in the form in the same way as before.
	if !form.Valid() {
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

	// Use the put() method to add a string value ("Snippet successfully created!") and the corresponding key ("flash") to the session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// w.Write([]byte("Create a new snippet..."))
	// Rediect the user to the relevant page for the snippet.
	// Update the redirect path to use the new clean URL format.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%v", id), http.StatusSeeOther)
}

// Create a new userSignForm struct
type signupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// Define a userSignupForm handelr func that displays the signup page.
func (app *application) userSignupForm(w http.ResponseWriter, r *http.Request) {
	templData := app.newTemplateData(r)
	templData.Form = signupForm{}

	app.render(w, http.StatusOK, "signup.html", templData)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	// Declare a zero-valued instance of our signupForm struct.
	var form signupForm

	// Parse the form data into the signupForm struct.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper funcs.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 16), "password", "This field must be at least 16 characters long")

	// If there are any errors, redisplay the signup form along with
	// a 422 status code.
	if !form.Valid() {
		templData := app.newTemplateData(r)
		templData.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", templData)
		return
	}

	// Create a new user record in the database, if the email exists
	// then add an error message to the form and re-display it.
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			templData := app.newTemplateData(r)
			templData.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", templData)
		} else {
			app.serverError(w, err)
		}

		return
	}

	// Else add a confirmation flash message to the session, confirming that their signup worked.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// Redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display an HTML form for logging in a user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login a user...")
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout a user...")
}
