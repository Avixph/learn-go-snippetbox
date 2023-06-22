package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Define a new comand-line flag with the name 'addr', a default value of
	// ":4000" and some short help text explaining what the flag controls.
	// The value of the flag will be stored in the addr var at runtime.
	addr := flag.String("addr", ":4000", "HTTP network address")

	// Importantly, we use the flag.Parse() func to parse the command-line
	// flag. This reads in the command-line flag value and assigns it to the
	// addr var. You need to call this *before* you use the addr var
	// otherwise it will always contain the default value of ":4000". If any
	// errs are encounted during parsing the application will be terminated.
	flag.Parse()

	// Use the http.NewServerMux() func to initialize a new servermux.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static"
	// directory. Note that the path given to the http.Dir func is relative
	// to the project directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.handle() func to register the file server as the handler
	// for all URL paths that start with "/static/". For matching paths, we
	// strip the "/static/" prefix before the request reaches the file
	// server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the home, snippetView and snippetCreate funcs as handlers
	// for the corrisponding URL patrerns with the servermux.
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// The value returned from the flag.Sring() func is a pointer to the
	// flag value, not the value itself. So we need to dereference the
	// pointer (i.e. prefix it with the * symbol) before using it. Note that
	// we're using the log.Printf() func to interpolate the address with the
	// log message.
	// Use the http.ListenAndServe() func to start a new web server. We pass
	// in two parameters: the TCP network address to listen on (ex:
	// (localhost::4000)) and the servermux we created. If http.
	// listenAndServe() returns an err we use the log.Fatal() func to log
	// the err message and exit. Note that any err returned by http.
	// listenAndServe() is always non-nill.
	log.Printf("Starting server on http://localhost%s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
