package wasm

import (
	"errors"
	"reflect"
	"syscall/js"
)

// ErrInvalidArgumentType is returned when a generated Go function wrapper receives invalid argument types from JS.
var ErrInvalidArgumentType = errors.New("invalid argument passed into Go function")

var errorType = reflect.TypeOf((*error)(nil)).Elem()

type goThrowable struct {
	Result js.Value `wasm:"result"`
	Error  js.Value `wasm:"error"`
}

// toJSFunc takes a reflect.Value of a Go function and converts it to a JS function that:
// Errors if the parameter types do not conform to the Go function signature,
// Throws an error if the last returned value is an error and is non-nil,
// Return an array if there's multiple non-error return values.
func toJSFunc(x reflect.Value) js.Value {
	funcType := x.Type()
	var hasError bool
	if funcType.NumOut() != 0 {
		hasError = funcType.Out(funcType.NumOut()-1) == errorType
	}

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		in, err := conformJSValueToType(funcType, this, args)
		if err != nil {
			return ToJSValue(goThrowable{
				Error: NewError(err),
			})
		}

		out := x.Call(in)

		if !hasError {
			return ToJSValue(goThrowable{
				Result: returnValue(out),
			})
		}

		lastParam := out[len(out)-1]
		if !lastParam.IsNil() {
			return ToJSValue(goThrowable{
				Error: NewError(lastParam.Interface().(error)),
			})
		}
		return ToJSValue(goThrowable{
			Result: returnValue(out[:len(out)-1]),
		})
	}).JSValue()
}

var jsValueType = reflect.TypeOf(js.Value{})

// conformJSValueToType attempts to convert the provided JS values to reflect.Values that match the
// types expected for the parameters of funcType.
func conformJSValueToType(funcType reflect.Type, this js.Value, values []js.Value) ([]reflect.Value, error) {
	if funcType.NumIn() == 0 {
		if len(values) != 0 {
			return nil, ErrInvalidArgumentType
		}
		return []reflect.Value{}, nil
	}

	if funcType.In(0) == jsValueType {
		// If the first parameter is a js.Value, it is assumed to be the value of `this`.
		values = append([]js.Value{this}, values...)
	}

	if funcType.IsVariadic() && funcType.NumIn()-1 > len(values) {
		return nil, ErrInvalidArgumentType
	}

	if !funcType.IsVariadic() && funcType.NumIn() != len(values) {
		return nil, ErrInvalidArgumentType
	}

	in := make([]reflect.Value, 0, len(values))
	for i, v := range values {
		paramType := funcType.In(i)
		x := reflect.Zero(paramType).Interface()
		err := FromJSValue(v, &x)
		if err != nil {
			return nil, err
		}

		in = append(in, reflect.ValueOf(x))
	}

	return in, nil
}

// returnValue wraps returned values by Go in a JS-friendly way.
// If there are no returned values, it returns undefined.
// If there is exactly one, it returns the JS equivalent.
// If there is more than one, it returns an array containing the JS equivalent of every returned value.
func returnValue(x []reflect.Value) js.Value {
	switch len(x) {
	case 0:
		return js.Undefined()
	case 1:
		return ToJSValue(x[0].Interface())
	}

	xInterface := make([]interface{}, 0, len(x))
	for _, v := range x {
		xInterface = append(xInterface, v.Interface())
	}

	return ToJSValue(xInterface)
}
