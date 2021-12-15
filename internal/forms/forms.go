package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

// Form  custom type struct hold the form values  embeds A url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initialize the form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		map[string][]string{},
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	resp := r.Form.Get(field)
	if resp == "" {
		return false
	}
	return true
}

// Required checks for required form fields
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Valid validate our form fields, returns true if there are no errors, otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// MinLength check for entered form field minimum length
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	submittedField := r.Form.Get(field)
	if len(submittedField) < length {
		f.Errors.Add(field, fmt.Sprint("This field must be at least %d characters", length))
		return false
	}
	return true

}

// IsEmail checks for valid emails
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email")
	}
}
