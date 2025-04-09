package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form nothing
type Form struct {
	url.Values
	Error errors
}

// Valid checks if there's error in form
func (f *Form) Valid() bool {
	return len(f.Error) == 0
}

// New creates new form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks for required fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Error.Add(field, "This field is required")
		}
	}
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Error.Add(field, "This field is required")
		return false
	}

	return true
}

// MinLength checks for string minimum length
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Error.Add(field, fmt.Sprintf("%s must be at least %d characters", field, length))
		return false
	}

	return true
}

// IsEmail checks for valid email address
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Error.Add(field, "Invalid email address")
	}
}
