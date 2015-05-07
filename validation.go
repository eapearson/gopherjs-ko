package ko

import "github.com/gopherjs/gopherjs/js"

type ValidatedObservable interface {
	IsValid() bool
}

func NewValidatedObservable(data interface{}) ValidatedObservable {
	return &Object{Global().Call("validatedObservable", data)}
}

func (ob *Object) IsValid() bool {
	return ob.Call("isValid").Bool()
}

type validation struct {
	*js.Object
}

func Validation() *validation {
	return &validation{Object: Global().Get("validation")}
}

func (v *validation) Init(config js.M) {
	v.Call("init", config)
}
