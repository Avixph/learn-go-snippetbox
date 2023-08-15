package main

import (
	"net/http"
	"testing"

	"github.com/Avixph/learn-go-snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// use the newTestApplication helper. Which for now just
	// contains a couple of mock loggers that discard anything written to them.
	app := newTestApplication(t)

	// Use the newTestServer helper to create a new test server, passing
	// in the value returned by our app.routes() method as the handler for the
	// server. This starts up an HTTPS server which listens on a randomly-chosen
	// port of our local machine for the duration of the test. Notice that we defer
	// a call to ts.Close() so that the server is shutdown when the test completes.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	// Check the value of the response status code and body using the same pattern
	// as before.
	assert.Equal(t, code, http.StatusOK)

	assert.Equal(t, string(body), "OK")
}
