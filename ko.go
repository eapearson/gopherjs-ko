// Package ko implements bindings to KnockoutJS.
// It also has bindings for the Knockout Validation library found on https://github.com/Knockout-Contrib/Knockout-Validation
// Using EnableSecureBinding make KnockoutJS works under CSP environments.
package ko

import (
	"github.com/Archs/js/dom"

	"github.com/gopherjs/gopherjs/js"
)

func ko() *js.Object {
	return js.Global.Get("ko")
}

type Observable struct {
	*js.Object
}

func NewObservable(args ...interface{}) *Observable {
	return &Observable{ko().Call("observable", args...)}
}

func (ob *Observable) Set(data interface{}) {
	ob.Object.Invoke(data)
}

func (ob *Observable) Get() *js.Object {
	return ob.Object.Invoke()
}

func (ob *Observable) Subscribe(fn func(*js.Object)) *Subscription {
	return &Subscription{
		Object: ob.Object.Call("subscribe", fn),
	}
}

func (ob *Observable) Extend(params js.M) *Observable {
	ob.Object.Call("extend", params)
	return ob
}

// The rateLimit extender, however, causes an observable to suppress and delay change notifications for a specified period of time. A rate-limited observable therefore updates dependencies asynchronously.
//
// The rateLimit extender can be applied to any type of observable, including observable arrays and computed observables. The main use cases for rate-limiting are:
//
// 		1. Making things respond after a certain delay
// 		2. Combining multiple changes into a single update
//
// when "notifyWhenChangesStop" is true change envent will be fired only after no change event detects anymore.
// "notifyWhenChangesStop" default is false, then it works under "notifyAtFixedRate" mode, at most one change in one timeframe.
func (ob *Observable) RateLimit(timeframeMS int, notifyWhenChangesStop ...bool) {
	method := "notifyAtFixedRate"
	if len(notifyWhenChangesStop) >= 1 && notifyWhenChangesStop[0] {
		method = "notifyWhenChangesStop"
	}
	ob.Extend(js.M{
		"rateLimit": js.M{
			"timeout": timeframeMS,
			"method":  method,
		},
	})
}

type ObservableArray struct {
	*Observable
}

func NewObservableArray(args ...interface{}) *ObservableArray {
	return &ObservableArray{&Observable{ko().Call("observableArray", args...)}}
}

func (ob *ObservableArray) IndexOf(data interface{}) int {
	return ob.Object.Call("indexOf", data).Int()
}

func (ob *ObservableArray) Pop() *js.Object {
	return ob.Object.Call("pop")
}

func (ob *ObservableArray) Unshift(data interface{}) {
	ob.Object.Call("unshift", data)
}

func (ob *ObservableArray) Shift() *js.Object {
	return ob.Object.Call("shift")
}

func (ob *ObservableArray) Reverse() {
	ob.Object.Call("reverse")
}

func (ob *ObservableArray) Sort() {
	ob.Object.Call("sort")
}

func (ob *ObservableArray) SortFunc(fn func(*js.Object, *js.Object)) {
	ob.Object.Call("sort", fn)
}

func (ob *ObservableArray) Splice(i, n int) *js.Object {
	return ob.Object.Call("splice", i, n)
}

func (ob *ObservableArray) RemoveAll(items ...interface{}) *js.Object {
	return ob.Object.Call("removeAll", items...)
}

func (ob *ObservableArray) Index(i int) *js.Object {
	return ob.Get().Index(i)
}

func (ob *ObservableArray) Length() int {
	return ob.Get().Length()
}

func (ob *ObservableArray) Push(data interface{}) {
	ob.Object.Call("push", data)
}

func (ob *ObservableArray) Remove(item interface{}) *js.Object {
	return ob.Object.Call("remove", item)
}

func (ob *ObservableArray) RemoveFunc(fn func(*js.Object) bool) *js.Object {
	return ob.Object.Call("remove", fn)
}

type Computed struct {
	*Observable
}

type WritableComputed struct {
	*Computed
}

func NewComputed(fn func() interface{}) *Computed {
	return &Computed{&Observable{ko().Call("computed", fn)}}
}

func NewWritableComputed(r func() interface{}, w func(interface{})) *WritableComputed {
	return &WritableComputed{
		&Computed{
			&Observable{
				ko().Call("computed", js.M{
					"read":  r,
					"write": w,
				}),
			},
		},
	}
}

func (ob *Computed) Dispose() {
	ob.Object.Call("dispose")
}

func (ob *Computed) Peek() *js.Object {
	return ob.Object.Call("peek")
}

type Subscription struct {
	*js.Object
}

func (s *Subscription) Dispose() {
	s.Object.Call("dispose")
}

type ComponentsFuncs struct {
	o *js.Object
}

func Components() *ComponentsFuncs {
	return &ComponentsFuncs{
		o: ko().Get("components"),
	}
}

func (co *ComponentsFuncs) Register(name string, params js.M) {
	co.o.Call("register", name, params)
}

// RegisterEx is an easy form to create KnockoutJS components
//  name is the component name
//  vmfunc is the ViewModel creator
//  template is the html tempalte for the component
//  cssRules would be directly embeded in the final html page, which can be ""
func (co *ComponentsFuncs) RegisterEx(name string, vmfunc func(params *js.Object) interface{}, template, cssRules string) {
	// embed the cssRules
	if cssRules != "" {
		style := dom.CreateElement("style")
		style.InnerHTML = cssRules
		dom.Body().AppendChild(style)
	}
	// register the component
	co.Register(name, js.M{
		"viewModel": vmfunc,
		"template":  template,
	})
}

func IsObservable(data interface{}) bool {
	return ko().Call("isObservable", data).Bool()
}

func IsComputed(data interface{}) bool {
	return ko().Call("isComputed", data).Bool()
}

func IsWritableObservable(data interface{}) bool {
	return ko().Call("isWritableObservable", data).Bool()
}

// RegisterURLTemplateLoader registers a new template loader which can be used to load
// template files from a webserver.
// To use it you need to pass a map with a `url` key as template argument to your component:
//   "template":  js.M{"url": "form.html"}
// This loader requires jQuery.
func RegisterURLTemplateLoader() {
	loader := func(name string, config *js.Object, callback func(*js.Object)) {
		url := config.Get("url")
		if url != nil && url != js.Undefined {
			// Some browsers are caching these requests too aggressively
			urlStr := url.String()
			urlStr += "?_=" + js.Global.Call("eval", `Date.now()`).String()

			js.Global.Get("jQuery").Call("get", urlStr, func(data *js.Object) {
				// We need an array of DOM nodes, not a string.
				// We can use the default loader to convert to the
				// required format.
				Components().o.Get("defaultLoader").Call("loadTemplate", name, data, callback)
			})
		} else {
			// Unrecognized config format. Let another loader handle it.
			callback(nil)
		}
	}

	Components().o.Get("loaders").Call("unshift", js.M{
		"loadTemplate": loader,
	})
}

func Unwrap(ob *js.Object) *js.Object {
	return ko().Call("unwrap", ob)
}

func ApplyBindings(args ...interface{}) {
	ko().Call("applyBindings", args...)
}
