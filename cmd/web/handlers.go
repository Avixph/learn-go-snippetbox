package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Define a home handler func that writes a byte slice containing
// "Hello from Snippetbox!" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If not, use
	// the htp.NotFound() func to send a 404 respond to the client.
	// Importantly, we then return from the handler. If we don't return, the
	// handler would keep executing and also write "Hello from Snippetbox!" message
	rup := r.URL.Path
	if rup != "/" {
		http.NotFound(w, r)
		return
	}

	// Initialize a slice containing the paths to the two templates. It's
	// important to note that the file containing our base template must be
	// the "first" file in the slice.
	tmplFiles := []string{
		"./ui/html/base.html",
		"./ui/html/components/nav.html",
		"./ui/html/pages/home.html",
		"./ui/html/components/footer.html",
	}

	// Use the template.ParseFiles() func to read the template files and
	// store the templates in a template set. If there's an err, we log a
	// detailed err message and use the http.Error() func to send a generic
	// 500 Interanl Server Err response to the user.
	ts, err := template.ParseFiles(tmplFiles...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// We then use the ExecuteTemplate() method to write the content of the
	// "base" template as the response body.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// Define a snippetView handler func
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extract the value of the id parameter from the query string and try to
	// convert it to an integer using the strconv.Atoi() func. If it can't be
	// converted to an integer, or it's value is less than 0, we return a 404
	// page not found response.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 0 {
		http.NotFound(w, r)
		return
	}

	// Use the fmt.Fprintf() func to interpolate the id value with our
	// response and write it to the http.ResponseWriter.
	fmt.Fprintf(w, "Displaying a specific snippet with ID# %d...", id)
}

// Define snippetCreate handler func
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	rm := r.Method
	if rm != http.MethodPost {
		// Use the Header().Set() method to add a 'Allow: POST' header to the
		// response header map. The first parameter is the header name, and
		// the second rarameter is the header value.
		w.Header().Set("Allow", http.MethodPost)

		// // If it's not, use the http.Error() func to send a 405 status code
		// and "Method Not Allowed!" string as the response body.
		http.Error(w, "Method not Allowed!", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
