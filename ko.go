package ko

import "github.com/gopherjs/gopherjs/js"

type Disposable interface {
	Dispose()
}

type Observable interface {
	Disposable

	Set(interface{})
	Get() *js.Object
	Subscribe(func(*js.Object)) Disposable
	Extend(js.M) Observable
}

type ObservableArray interface {
	Observable

	Index(int) *js.Object
	Length() int
	Push(interface{})
	Remove(interface{}) *js.Object
	RemoveFunc(func(*js.Object) bool) *js.Object
}

type Computed interface {
	Observable
}

type ValidatedObservable interface {
	IsValid() bool
}

type Object struct {
	*js.Object
}

func (ob *Object) Dispose() {
	ob.Call("dispose")
}

func (ob *Object) Set(data interface{}) {
	ob.Invoke(data)
}

func (ob *Object) Get() *js.Object {
	return ob.Invoke()
}

func (ob *Object) Subscribe(fn func(*js.Object)) Disposable {
	return &Object{ob.Call("subscribe", fn)}
}

func (ob *Object) Extend(params js.M) Observable {
	ob.Call("extend", params)
	return ob
}

func (ob *Object) Index(i int) *js.Object {
	return ob.Get().Index(i)
}

func (ob *Object) Length() int {
	return ob.Get().Length()
}

func (ob *Object) IndexOf(data interface{}) int {
	return ob.Call("indexOf", data).Int()
}

func (ob *Object) Push(data interface{}) {
	ob.Call("push", data)
}

func (ob *Object) Pop() *js.Object {
	return ob.Call("pop")
}

func (ob *Object) Unshift(data interface{}) {
	ob.Call("unshift", data)
}

func (ob *Object) Shift() *js.Object {
	return ob.Call("shift")
}

func (ob *Object) Reverse() {
	ob.Call("reverse")
}

func (ob *Object) Sort() {
	ob.Call("sort")
}

func (ob *Object) SortFunc(fn func(*js.Object, *js.Object)) {
	ob.Call("sort", fn)
}

func (ob *Object) Splice(i, n int) *js.Object {
	return ob.Call("splice", i, n)
}

func (ob *Object) Remove(item interface{}) *js.Object {
	return ob.Call("remove", item)
}

func (ob *Object) RemoveFunc(fn func(*js.Object) bool) *js.Object {
	return ob.Call("remove", fn)
}

func (ob *Object) RemoveAll(items ...interface{}) *js.Object {
	return ob.Call("removeAll", items...)
}

func (ob *Object) IsValid() bool {
	return ob.Call("isValid").Bool()
}

type components struct {
	*js.Object
}

func Components() *components {
	return &components{
		Object: Global().Get("components"),
	}
}

func (co *components) Register(name string, params js.M) {
	co.Call("register", name, params)
}

func Global() *js.Object {
	return js.Global.Get("ko")
}

func NewObservable(data interface{}) Observable {
	return &Object{Global().Call("observable", data)}
}

func NewObservableArray(data interface{}) ObservableArray {
	return &Object{Global().Call("observableArray", data)}
}

func NewValidatedObservable(data interface{}) ValidatedObservable {
	return &Object{Global().Call("validatedObservable", data)}
}

func NewComputed(fn func() interface{}) Computed {
	return &Object{Global().Call("computed", fn)}
}

func RegisterTemplateURLLoader() {
	loader := func(name string, config *js.Object, callback func(*js.Object)) {
		url := config.Get("url")
		if url != nil && url != js.Undefined {
			js.Global.Get("jQuery").Call("get", url, func(data *js.Object) {
				// We need an array of DOM nodes, not a string.
				// We can use the default loader to convert to the
				// required format.
				Components().Get("defaultLoader").Call("loadTemplate", name, data, callback)
			})
		} else {
			// Unrecognized config format. Let another loader handle it.
			callback(nil)
		}
	}

	Components().Get("loaders").Call("unshift", js.M{
		"loadTemplate": loader,
	})
}

func ApplyBindings(model interface{}) {
	Global().Call("applyBindings", model)
}
