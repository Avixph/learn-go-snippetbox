package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Define an application struct to hold the app-wide dependencies for the
// web app. For now we'll only include feilds for the two custom loggers,
// but  we'll add more to it as the build progresses.
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

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

	// Use log.New() to create a logger for writting info messages. This
	// takes three parameters: the destination to wrie the logs to (os.
	// Stdout), a string prefix for message (INFO followed by tab), and
	// flags to indicate what additional info ito include (local data and
	// time). Note thatthe flags ae joined using the bitwise OR operator |.
	// Create a logger for writting error messages in the same way, but use
	// stderr as the destination and use the log.Lshortfile to include the
	// relevant file name and line number.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of our application struct, containing the
	// dependencies
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// Use the http.NewServerMux() func to initialize a new servermux.
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static"
	// directory. Note that the path given to the http.Dir func is relative
	// to the project directory root.
	// Use the mux.handle() func to register the file server as the handler
	// for all URL paths that start with "/static/". For matching paths, we
	// strip the "/static/" prefix before the request reaches the file
	// server.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the home, snippetView and snippetCreate funcs as handlers
	// for the corrisponding URL patrerns with the servermux.
	// Swap the route declearations to use the application struct's methods
	// as the handler func.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Initalize a new http.Server struct. We set the Addr and Handler
	// fields so that the server uses the same network address and routes as
	// before, and set the ErrorLog fielfd so that ther server now uses the
	// custom errorLog logger in the event of and problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// The value returned from the flag.Sring() func is a pointer to the
	// flag value, not the value itself. So we need to dereference the
	// pointer (i.e. prefix it with the * symbol) before using it. Note that
	// we're using the infoLog.Printf() func to interpolate the address with
	// the log message.
	// Use the http.ListenAndServe() func on the http.Server() struct to
	// start a new web server. We pass in two parameters: the TCP network
	// address to listen on (ex:(localhost::4000)) and the servermux we
	// created. If http. listenAndServe() returns an err we use the errorLog.
	// Fatal() func to log the err message and exit. Note that any err
	// returned by http. listenAndServe() is always non-nill.
	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
