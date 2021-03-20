package wasm

import (
	"fmt"
	"reflect"
	"syscall/js"
	"unsafe"
)

// ToJSValue converts a given Go value into its equivalent JS form.
//
// One special case is that complex numbers (complex64 and complex128) are converted into objects with a real and imag
// property holding a number each.
//
// A function is converted into a JS function where the function returns an error if the provided arguments do not conform
// to the Go equivalent but otherwise calls the Go function.
//
// The "this" argument of a function is always passed to the Go function if its first parameter is of type js.Value.
// Otherwise, it is simply ignored.
//
// If the last return value of a function is an error, it will be thrown in JS if it's non-nil.
// If the function returns multiple non-error values, it is converted to an array when returning to JS.
//
// It panics when a channel or a map with keys other than string and integers are passed in.
func ToJSValue(x interface{}) js.Value {
	if x == nil {
		return js.Null()
	}

	// Fast path for basic types that do not require reflection.
	switch x := x.(type) {
	case js.Value:
		return x
	case js.Wrapper:
		return x.JSValue()
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr,
		unsafe.Pointer, float32, float64, string:
		return js.ValueOf(x)
	case complex64:
		return js.ValueOf(map[string]interface{}{
			"real": real(x),
			"imag": imag(x),
		})
	case complex128:
		return js.ValueOf(map[string]interface{}{
			"real": real(x),
			"imag": imag(x),
		})
	}

	value := reflect.ValueOf(x)

	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
		if !value.IsValid() {
			return js.Undefined()
		}
	}

	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		return toJSArray(value)
	case reflect.Func:
		return toJSFunc(value)
	case reflect.Map:
		return mapToJSObject(value)
	case reflect.Struct:
		return structToJSObject(value)
	default:
		panic(fmt.Sprintf("cannot convert %T to a JS value", x))
	}
}

func toJSArray(x reflect.Value) js.Value {
	arrayConstructor, err := Global().Get("Array")
	if err != nil {
		panic("Array constructor not found")
	}

	array := arrayConstructor.New()
	for i := 0; i < x.Len(); i++ {
		array.SetIndex(i, ToJSValue(x.Index(i).Interface()))
	}

	return array
}

func mapToJSObject(x reflect.Value) js.Value {
	objectConstructor, err := Global().Get("Object")
	if err != nil {
		panic("Object constructor not found")
	}

	obj := objectConstructor.New()
	iter := x.MapRange()
	for {
		if !iter.Next() {
			break
		}

		key := iter.Key()
		value := iter.Value().Interface()
		switch key := key.Interface().(type) {
		case int:
			obj.SetIndex(key, ToJSValue(value))
		case int8:
			obj.SetIndex(int(key), ToJSValue(value))
		case int16:
			obj.SetIndex(int(key), ToJSValue(value))
		case int32:
			obj.SetIndex(int(key), ToJSValue(value))
		case int64:
			obj.SetIndex(int(key), ToJSValue(value))
		case uint:
			obj.SetIndex(int(key), ToJSValue(value))
		case uint8:
			obj.SetIndex(int(key), ToJSValue(value))
		case uint16:
			obj.SetIndex(int(key), ToJSValue(value))
		case uint32:
			obj.SetIndex(int(key), ToJSValue(value))
		case uint64:
			obj.SetIndex(int(key), ToJSValue(value))
		case uintptr:
			obj.SetIndex(int(key), ToJSValue(value))
		case string:
			obj.Set(key, ToJSValue(value))
		default:
			panic(fmt.Sprintf("cannot convert %T into a JS value as its key is not a string or an integer",
				x.Interface()))
		}
	}

	return obj
}

func structToJSObject(x reflect.Value) js.Value {
	objectConstructor, err := Global().Get("Object")
	if err != nil {
		panic("Object constructor not found")
	}

	obj := objectConstructor.New()

	structType := x.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name := field.Name
		if tagName, ok := field.Tag.Lookup("wasm"); ok {
			name = tagName
		}

		obj.Set(name, ToJSValue(x.Field(i).Interface()))
	}

	return obj
}
