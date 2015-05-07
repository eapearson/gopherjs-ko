package components

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mibitzi/gopherjs-ko"
)

// SetupHome returns a component configuration for the home component.
func SetupHome() js.M {
	return js.M{
		"viewModel": NewHome,
		"template":  js.M{"url": "home.html"},
	}
}

// NewHome returns a new viewModel for the home component.
func NewHome(params *js.Object) js.M {
	return js.M{
		"name": ko.NewObservable(""),
	}
}
