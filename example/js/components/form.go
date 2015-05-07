package components

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mibitzi/gopherjs-ko"
)

// Form contains observables used to store form data.
type Form struct {
	email ko.Observable
	phone ko.Observable
	form  ko.ValidatedObservable
}

// SetupForm returns a component configuration for the form component.
func SetupForm() js.M {
	return js.M{
		"viewModel": NewForm,
		"template":  js.M{"url": "form.html"},
	}
}

// NewForm returns a new viewModel for the form component.
func NewForm(params *js.Object) js.M {
	form := &Form{
		// We use knockout validation to validate email and phone fields.
		// https://github.com/Knockout-Contrib/Knockout-Validation

		email: ko.NewObservable("").Extend(js.M{
			"required": true,
			"email":    true,
		}),

		phone: ko.NewObservable("").Extend(js.M{
			"required": true,
		}).Extend(js.M{
			"minLength": 5,
		}).Extend(js.M{
			"pattern": js.M{
				"message": "This doesn't look like a phone number to me",
				"params":  "^[0-9 -]+$",
			},
		}),
	}

	// We use a ValidatedObservable to check if the form is valid.
	form.form = ko.NewValidatedObservable(js.S{form.email, form.phone})

	return js.M{
		"email": form.email,
		"phone": form.phone,
		"save":  form.save,
	}
}

func (form *Form) save() {
	if !form.form.IsValid() {
		js.Global.Call("alert", "Form is not valid")
	} else {
		js.Global.Call("alert", "Form is valid")
	}
}
