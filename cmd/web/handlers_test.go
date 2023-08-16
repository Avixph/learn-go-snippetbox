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

func TestSnippetView(t *testing.T) {
	// Create a new instance of the application struct which uses the mocked
	// dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running end-to-end tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Set up table-driven tests to check the responses sent by the app for
	// different URLs.
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/3fa80338-89b5-407e-b294-c3ac68238070",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Integer ID",
			urlPath:  "/snippet/view/52898447124421815065728728672829",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-98030664472603982489384866019669",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}

}
