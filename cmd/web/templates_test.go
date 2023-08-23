package main

import (
	"testing"
	"time"

	"github.com/Avixph/learn-go-snippetbox/internal/assert"
)

func TestHuanDate(t *testing.T) {
	// Create a slice of annonymous structs containg the test case name, input to our
	// humanDate func (the tm field which is a new time.Time object).
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2022 at 09:15",
		},
	}

	// Loop over the test cases.
	for _, tt := range tests {
		// Use the t.Run() func to run a sub-test for each test case. The first parameter to this is the name of the test (which is used to identify yhe sub-test in any log output) and the second parameter is an anonymous function containing the actual test for each test case.
		t.Run(tt.name, func(t *testing.T) {
			// Check that the ouput from the humaenDate func is in our expected format. If
			// it isn't, then use the t.Errof() func to indicate that the test failed and
			// log the expected/ actual values.
			hd := humanDate(tt.tm)

			// Use the assert.Equal helper to compare the expected and actual values.
			assert.Equal(t, hd, tt.want)
		})
	}
}
