// This example shows how to implement a basic single page application using gopherjs-ko
// in combination with page.js
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mibitzi/gopherjs-ko"
	"github.com/mibitzi/gopherjs-ko/example/js/components"
)

func main() {
	// Each page is one component (which can contain other components).
	ko.Components().Register("home-page", components.SetupHome())
	ko.Components().Register("form-page", components.SetupForm())

	// Register a new template loader which can load *.html files from our server.
	// This loader requires jQuery.
	ko.RegisterURLTemplateLoader()

	// We need to configure knockout validation to use bootstrap syntax for messages.
	ko.Validation().Init(js.M{
		"errorMessageClass": "help-block",
		"errorElementClass": "has-error",
	})

	router := NewRouter()

	ko.ApplyBindings(js.M{
		// These bindings will be used to display the current page and to set the .active class
		// in the navbar.
		"route": router.Current,
		"path":  router.Path,
	})
}
