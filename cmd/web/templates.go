package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/Avixph/learn-go-snippetbox/ui"
)

// Define a templateData type to act as the holding structure for any
// dynamic data that we want to pass to our HTML templates. At the moment it
// only contains one feild, but we'll add more to it as the build progresses.
// Add a Form field with the type "any" a Flash field, a IsAuthenticated field,
// and a CSRFToken field to the templateData struct.
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

// Create a humanDate func that returns a nicely formatted string
// representation of time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable.
// This is essentially a string-keyed map that acts as a lookup between the
// names of our custom template funcs and the funcs themselves.
var templFunctions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the fs.Glob() func to get a slice of all filepaths in the ui.Files
	// embeded filesystem that match 'html/pages/*.html'. This will essentially
	// give us a slice of all the 'page' templates of the app.
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through each of the page filepaths.
	for _, page := range pages {
		// Extract the file name (ex: 'home.html') from the full filepath and
		// assign it to the name variable.
		name := filepath.Base(page)

		// Create a slice containing the filepath patterns for the templates we want
		// to parse.
		patterns := []string{
			"html/base.html",
			"html/components/*.html",
			page,
		}

		// The template.FuncMap must be registered with the template set before calling 
		// the ParseFS() method. This means we have to use template.New() to create an 
		// empty template set, use the Funcs() method to register the template.FuncMap, 
		// and then parse the tempalte file.
		ts, err := template.New(name).Funcs(templFunctions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
