package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Avixph/learn-go-snippetbox/internal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

// Define an application struct to hold the app-wide dependencies for the
// web app. For now we'll only include feilds for the two custom loggers.
// Add a snippets field to the application struct. This will allow us to
// make the SnippetModel object available to our handlers.
// Addd a templateCache feild, formDecoder field, a sessionManager field,
// a users field and a debug field to the application struct.
type application struct {
	debug          bool
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
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

	// Define a new debug flag with a defualt value of false.
	debug := flag.Bool("debug", false, "Enable debug mode")

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
	// Make sure that the Secure attribute is set on our session coockies.
	// Setting ths means that the cookie will only be sent by a user's web
	// browser when HTTPS connection is being used (and wo't be sent over
	// unsecure HTTP connections).
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initialize a new instance of our application struct, containing the
	// dependencies.
	// Initialize a models.SnippetModel instance and add it to the
	// application dependencies.
	// Add a templateCache, a formDecoder, a sessionManager, a models.
	// UserModel, and a debug to the application dependencies.
	app := &application{
		debug:          *debug,
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve prefernece value, so that the only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Initalize a new http.Server struct. We set the Addr and Handler
	// fields so that the server uses the same network address and routes as
	// before, and set the ErrorLog fielfd so that ther server now uses the
	// custom errorLog logger in the event of and problems.
	// Call the new app.routes() method to get the servermux containing our
	// routes.
	// Set the server's TLSConfig field to use the tlsConfig variable.
	// Add Idle, Read Write timeouts to the server.
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// The value returned from the flag.Sring() func is a pointer to the
	// flag value, not the value itself. So we need to dereference the
	// pointer (i.e. prefix it with the * symbol) before using it. Note that
	// we're using the infoLog.Printf() func to interpolate the address with
	// the log message.
	// Use the http.ListenAndServeTLS() func on the http.Server() struct to
	// start a new web server (passing in the paths to the TLS certificate and
	//corresponding private key). We pass in two parameters: the TCP network
	// address to listen on (ex:(localhost::4000)) and the servermux we
	// created. If http. listenAndServe() returns an err we use the errorLog.
	// Fatal() func to log the err message and exit. Note that any err
	// returned by http. listenAndServe() is always non-nill.
	// Because the err var is already declared above, we need to use the
	// assignment operator "=" here, instead of ":=" 'declare and assigng'
	infoLog.Printf("Starting server on https://localhost%s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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

// CREATE DATABASE test_snippetbox WITH ENCODING 'UTF8' LC_COLLATE='en_US.UTF-8' LC_CTYPE='en_US.UTF-8' TEMPLATE=template0;

// CREATE USER test_web WITH PASSWORD 'learn-go-snippetbox';
// // GRANT CREATE, SELECT, INSERT, UPDATE, DELETE ON DATABASE test_snippetbox TO test_web;
// GRANT ALL ON DATABASE test_snippetbox TO test_web;
// REVOKE TEMPORARY ON DATABASE test_snippetbox FROM test_web;

// INSERT INTO snippets (title, content, created_on, updated_on, expires_on) VALUES (
//     'First autumn morning',
//     'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
//     (now() at time zone 'utc'),
//     (now() at time zone 'utc'),
//                 (now() at time zone 'utc' + interval '7 day')
// );
