package components

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mibitzi/gopherjs-ko"
)

type Home struct {
	data ko.Observable
	text ko.Observable
}

// SetupHome returns a component configuration for the home component.
func SetupHome() js.M {
	return js.M{
		"viewModel": NewHome,
		"template":  js.M{"url": "home.html"},
	}
}

// NewHome returns a new viewModel for the home component.
func NewHome(params *js.Object) js.M {
	home := &Home{
		data: ko.NewObservable(""),
		text: ko.NewObservable(""),
	}

	return js.M{
		"data":      home.data,
		"text":      home.text,
		"resetText": home.resetText,
		"fetchData": home.fetchData,
	}
}

func (home *Home) resetText() {
	home.text.Set("")
}

func (home *Home) fetchData() {
	js.Global.Get("jQuery").Call("get", "data.json", func(data *js.Object, status *js.Object, xhr *js.Object) {
		home.data.Set(xhr.Get("responseText"))
	})
}
