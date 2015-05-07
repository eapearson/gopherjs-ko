package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/mibitzi/gopherjs-ko"
	"honnef.co/go/js/dom"
)

type Router struct {
	router  *js.Object
	Current ko.Observable
	Active  ko.Computed
}

func NewRouter() *Router {
	r := &Router{
		Current: ko.NewObservable(js.M{}),
	}

	r.Active = ko.NewComputed(func() interface{} {
		return r.Current.Get().Get("path")
	})

	r.init()

	return r
}

func (r *Router) init() {
	// We use pagejs to do our routing.
	page := js.Global.Get("page")

	// Handle the home page
	page.Invoke("/", func(ctx *js.Object) {
		r.Set("home-page", ctx)
	})

	// Handle the form page
	page.Invoke("/form", func(ctx *js.Object) {
		r.Set("form-page", ctx)
	})

	// Redirect everything else to /
	page.Call("redirect", "*", "/")

	// Apply our paths and config
	page.Invoke(js.M{})
}

// Set changes the current page and sets its parameters.
func (r *Router) Set(p string, ctx *js.Object) {
	// Whenever we change the page we should also reset the scroll position.
	dom.GetWindow().ScrollTo(0, 0)

	r.Current.Set(js.M{
		"page":   p,
		"params": ctx.Get("params"),
		"path":   ctx.Get("path"),
	})
}
