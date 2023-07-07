package main

import (
	"html/template"
	"path/filepath"

	"github.com/Avixph/learn-go-snippetbox/internal/models"
)

// Define a templateData type to act as the holding structure for any
// dynamic data that we want to pass to our HTML templates. At the moment it
// only contains one feild, but we'll add more to it as the build progresses.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() func to get a slice of all filepaths that match
	// the pattern ".ui/html/pages/*.html". This will essentially give us a
	// slice of all the filepaths of the app 'page' templates.
	// Ex: [ui/html/pages/home.html ui/html/pages/view.html]
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// Loop through each of the page filepaths.
	for _, page := range pages {
		// Extract the file name (ex: 'home.html') from the full filepath and
		// assign it to the name variable.
		name := filepath.Base(page)

		// Parse the base template file into a template set.
		// ts, err := template.ParseFiles("./ui/html/base.html")
		// if err != nil {
		// 	return nil, err
		// }

		// Call ParseGlob() *on this template set* to add any componemts.
		// ts, err = template.ParseGlob("./ui/html/components/*.html")
		// if err != nil {
		// 	return nil, err
		// }

		// Call ParseFile() *on this template set* to add the page template.
		// ts, err = template.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }

		// Create a slice containing the filepaths for our base template, any
		// partials and the page.
		tmplFiles := []string{
			"./ui/html/base.html",
			"./ui/html/components/nav.html",
			page,
			"./ui/html/components/footer.html",
		}

		// Parse the files into a template set.
		ts, err := template.ParseFiles(tmplFiles...)
		if err != nil {
			return nil, err
		}

		// Add the tempalte set to the map, using the name of the page (ex:
		// 'home.html') as the key.
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
