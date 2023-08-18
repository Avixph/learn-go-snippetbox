package main

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/Avixph/learn-go-snippetbox/internal/models/mocks"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

// Define a regexp (regular expression) which captures the CSRF token value
// from the HTML for our user signup page.
var csrfTokenRX = regexp.MustCompile(`<input type="hidden" name="csrf_token" value="(.+)" />`)

// Create a extractCSRFToken helper which returns an html unescape string.
func extractCSRFToken(t *testing.T, body string) string {
	// Use the FindStringSubmatch method to extract the token from the HTML body.
	// Note that this returns an array with that entire matched pattern in the first
	// position, and the values of any captured data in the subsequent positions.
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(string(matches[1]))
}

// Create a newTestApplication helper which returns an instance of our
// application struct  containing mocked dependencies.
func newTestApplication(t *testing.T) *application {
	// Create an insance of the template cache and a form decoder.
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	// Create a session manager instance using the same settings as production,
	// except that we *don't* set a Store for the session manager. If no store is
	// set, the SCS package will default to using a transient in-memory store,
	// which is ideal for testing purposes.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLog:       log.New(io.Discard, "", 0),
		infoLog:        log.New(io.Discard, "", 0),
		snippets:       &mocks.SnippetModel{},
		users:          &mocks.UserModel{},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// Create a newTestServer helper which initializes and returns a new instance
// of our custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize the test server
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client which when used will store and send any response cookies with subsequent requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the test server client by setting a custom CheckRedirect func. THis func will be called whenever a 3xx response is received
	// by the client, and by always returning a http.ErrUseLastResponse error it
	// forces the client to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// Implement a get() method on our custom testServer type. This makes a GET
// request to a given url path using the test server client, and returns the
// response status code, headers and body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	// Check that the response body written by the ping equals "OK".
	// The network address that the test server is listening on and is contained
	// in the ts.URL field. We can use this along with the ts.Client().Get() method
	// to make a Get/ping request against the test server. This returns a http.Response
	// struct containing the rsponse.
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

// Implement a postForm method for sending POST requests to the test server. The
// final parameter to this method is a url.Values object which can contain any form data that sent in the request body.
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	// Read the response body from the test server.
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	// Return the response status, headers and body.
	return rs.StatusCode, rs.Header, string(body)
}
