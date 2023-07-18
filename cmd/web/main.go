package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	// Import the models package that we created.
	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

// Define an application struct to hold the app-wide dependencies for the
// web app. For now we'll only include feilds for the two custom loggers.
// Add a snippets field to the application struct. This will allow us to make
// the SnippetModel object available to our handlers.
// Addd a templateCache feild, formDecoder field, and a sessionManager field
// to the application struct.
type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	// Define a new comand-line flag with the name 'addr', a default value of
	// ":4000" and some short help text explaining what the flag controls.
	// The value of the flag will be stored in the addr var at runtime.
	// Use the getEnvVariables() helper
	addr := flag.String("addr", getEnvVariables("LOCAL_PORT"), "HTTP network address")

	// Define a new command-line flag for the PostgreSQL DSN string.
	dsn := flag.String("dsn", getEnvVariables("DATABASE_URL"), "PostgresSQL data source name")

	// Importantly, we use the flag.Parse() func to parse the command-line
	// flag. This reads in the command-line flag value and assigns it to the
	// addr var. We need to call this *before* using the addr var
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

	// To keep the main() func tidy the code for creating a connection pool was
	// placed into a seperate openDB() func. WE pass opewnDB() the DSN from the
	// command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// We also defer a call to the db.Close(), so that the connection pool is
	// closed before the main() func exits.
	defer db.Close()

	// Initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initialize a decoder instance.
	formDecoder := form.NewDecoder()

	// Initialize a new session manager with scs.New() funct. Then we configure
	// it touse oour PostgeSQL database as the session store, and set a lifetime
	// of 12 hours (so that sessions automatically expire after 12 hours of
	// creation.)
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	// Initialize a models.SnippetModel instance and add it to the application
	// dependencies.
	// Add a templateCache, a formDecoder, and a sessionManager to the
	// application dependencies.
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initalize a new http.Server struct. We set the Addr and Handler
	// fields so that the server uses the same network address and routes as
	// before, and set the ErrorLog fielfd so that ther server now uses the
	// custom errorLog logger in the event of and problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		// Call the new app.routes() method to get the servermux containing our routes.
		Handler: app.routes(),
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
	// Because the err var is already declared above, we need to use the
	// assignment operator "=" here, instead of ":=" 'declare and assigng'
	infoLog.Printf("Starting server on http://localhost%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// The openDB() func wraps sql.Open() and returns a sql.DB connection pool for
// the given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Use the db.Ping() method to create a connection and check for any errors.
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
