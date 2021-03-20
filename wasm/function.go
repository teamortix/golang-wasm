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

func conformJSValueToType(funcType reflect.Type, this js.Value, values []js.Value) ([]reflect.Value, error) {
	if funcType.NumIn() == 0 {
		if len(values) != 0 {
			return nil, ErrInvalidArgumentType
		}
		return []reflect.Value{}, nil
	}

	if funcType.In(0) == jsValueType {
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
