package main

import (
	"log"
	"net/http"
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

	w.Write([]byte("Hello from Snippetbox!"))
}

// Define a snippetView handler func
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Displaying a specific snippet..."))
}

// Define snippetCreate handler func
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// Use the http.NewServerMux() func to initialize a new servermux, then
	// register the home func as a handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	// Register the snippetView and snippetCreate handler funcs and
	// corrisponding URL patrerns with the servermux.
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Use the http.ListenAndServe() func to start a new web server. We pass
	// in two parameters: the TCP network address to listen on (ex:(localhost::4000))
	// and the servermux we created. If http.listenAndServe() returns an err
	// we use the log.Fatal() func to log the err message and exit. Note that
	// any err returned by http.listenAndServe() is always non-nill.
	log.Print("Starting server on localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
