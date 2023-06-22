package main

import (
	"log"
	"net/http"
)

// Define a home handler func that writes a byte slice containing
// "Hello from Snippetbox!" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox!"))
}

func main() {
	// Use tghe http.NewServerMux() func to initialize a new servermux, then
	// register the home func as a handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Use the http.ListenAndServe() func to start a new web server. We pass
	// in two parameters: the TCP network address to listen on (ex:(localhost::4000))
	// and the servermux we created. If http.listenAndServe() returns an err
	// we use the log.Fatal() func to log the err message and exit. Note that
	// any err returned by http.listenAndServe() is always non-nill.
	log.Print("Starting server on localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
