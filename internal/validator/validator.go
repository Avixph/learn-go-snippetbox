package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// expression pattern for  checking the format of an email
// address. This returns a pointer to a 'compiled' regexp.
// Regexp type, or panic if there's an error. Parsing this
// pattern once at startup and storing the compiled *regexp.
//
//	in a variable is more performant than re-parsin the
//
// pattern each time we need it.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Define a new validator type which contains a map of
// validation errors for our form fields.
// Add a new NonFieldErrors []string field to the struct,
// which we will use to hold any validation errors which
// are not related to a specific form field.
type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// The Valid() method returns true if the FieldErrors map and
// NonFieldErrors slice don't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// The AddNonFieldError() method adds error messages to
// the new NonFieldErrors slice.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// The AddFieldError() method adds an error message to the
// FieldErrors map (so long as no entry already exists for
// the given key).
func (v *Validator) AddFieldError(key, message string) {
	// Note: We need to initialize the map first, if it isn't
	// already initialized.
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// The CheckField() method adds an error message to the
// FieldErrors map only if a validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// The NotBlanck() func returns true if a value is not an
// empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// The MaxChars() func returns true if a value contains no
// more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// The generic PermittedValue() func returns true only if the value of type 
// T equals one of the variadic permitted values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// The MinChars() func returns true if a
// value contains at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// The Matches() func returns true if a value
// matches a provided compiled regular
// expression pattern.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
