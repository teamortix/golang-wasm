package wasm

import (
	"errors"
	"fmt"
	"reflect"
	"syscall/js"
)

// ErrMultipleReturnValue is an error where a JS function is attempted to be unmarshalled into a Go function with
// multiple return values.
var ErrMultipleReturnValue = errors.New("a JS function can only return one value")

// InvalidFromJSValueError is an error where an invalid argument is passed to FromJSValue.
// The argument to Unmarshal must be a non-nil pointer.
type InvalidFromJSValueError struct {
	Type reflect.Type
}

// Error implements error.
func (e InvalidFromJSValueError) Error() string {
	return "invalid argument passed to FromJSValue. Got type " + e.Type.String()
}

// InvalidTypeError is an error where the JS value cannot be unmarshalled into the provided Go type.
type InvalidTypeError struct {
	JSType js.Type
	GoType reflect.Type
}

// Error implements error.
func (e InvalidTypeError) Error() string {
	return "invalid unmarshalling: cannot unmarshal " + e.JSType.String() + " into " + e.GoType.String()
}

// InvalidArrayError is an error where the JS's array length do not match Go's array length.
type InvalidArrayError struct {
	Expected int
	Actual   int
}

// Error implements error.
func (e InvalidArrayError) Error() string {
	return fmt.Sprintf(
		"invalid unmarshalling: expected array of length %d to match Go array but got JS array of length %d",
		e.Expected, e.Actual,
	)
}

// Decoder is an interface which manually decodes js.Value on its own.
// It overrides in FromJSValue.
type Decoder interface {
	FromJSValue(js.Value) error
}

// FromJSValue converts a given js.Value to the Go equivalent.
// The new value of 'out' is undefined if FromJSValue returns an error.
//
// When a JS function is unmarshalled into a Go function with only one return value, the returned JS value is casted
// into the type of the return value. If the conversion fails, the function call panics.
//
// When a JS function is unmarshalled into a Go function with two return values, the second one being error, the
// conversion error is returned instead.
func FromJSValue(x js.Value, out interface{}) error {
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return &InvalidFromJSValueError{reflect.TypeOf(v)}
	}

	return decodeValue(x, v.Elem())
}

// decodeValue decodes the provided js.Value into the provided reflect.Value.
func decodeValue(x js.Value, v reflect.Value) error {
	// If we have undefined or null, we need to be able to set to the pointer itself.
	// All code beyond this point are pointer-unaware so we handle undefined or null first.
	if x.Type() == js.TypeUndefined || x.Type() == js.TypeNull {
		return decodeNothing(v)
	}

	// Implementations of Decoder are probably on pointer so do it before pointer code.
	if d, ok := v.Addr().Interface().(Decoder); ok {
		return d.FromJSValue(x)
	}

	// Make sure everything is initialized and indirect it.
	// This prevents other decode functions from having to handle pointers.
	if v.Kind() == reflect.Ptr {
		initializePointerIfNil(v)
		v = reflect.Indirect(v)
	}

	if v.Kind() == reflect.Interface && v.NumMethod() == 0 {
		// It's a interface{} so we just create the easiest Go representation we can in createInterface.
		res := createInterface(x)
		if res != nil {
			v.Set(reflect.ValueOf(res))
		}
		return nil
	}

	// Directly set v if it's a js.Value.
	if _, ok := v.Interface().(js.Value); ok {
		v.Set(reflect.ValueOf(x))
		return nil
	}

	// Go the reflection route.
	switch x.Type() {
	case js.TypeBoolean:
		return decodeBoolean(x, v)
	case js.TypeNumber:
		return decodeNumber(x, v)
	case js.TypeString:
		return decodeString(x, v)
	case js.TypeSymbol:
		return decodeSymbol(x, v)
	case js.TypeObject:
		if isArray(x) {
			return decodeArray(x, v)
		}
		return decodeObject(x, v)
	case js.TypeFunction:
		return decodeFunction(x, v)
	default:
		panic("unknown JS type: " + x.Type().String())
	}
}

// decodeNothing decodes an undefined or a null into the provided reflect.Value.
func decodeNothing(v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return InvalidTypeError{js.TypeNull, v.Type()}
	}
	v.Set(reflect.Zero(v.Type()))
	return nil
}

// decodeBoolean decodes a bool into the provided reflect.Value.
func decodeBoolean(x js.Value, v reflect.Value) error {
	if v.Kind() != reflect.Bool {
		return InvalidTypeError{js.TypeBoolean, v.Type()}
	}
	v.SetBool(x.Bool())
	return nil
}

// decodeNumber decodes a JS number into the provided reflect.Value, truncating as necessary.
func decodeNumber(x js.Value, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(int64(x.Float()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(uint64(x.Float()))
	case reflect.Float32, reflect.Float64:
		v.SetFloat(x.Float())
	default:
		return InvalidTypeError{js.TypeNumber, v.Type()}
	}
	return nil
}

// decodeString decodes a JS string into the provided reflect.Value.
func decodeString(x js.Value, v reflect.Value) error {
	if v.Kind() != reflect.String {
		return InvalidTypeError{js.TypeString, v.Type()}
	}
	v.SetString(x.String())
	return nil
}

// decodeSymbol decodes a JS symbol into the provided reflect.Value.
func decodeSymbol(x js.Value, v reflect.Value) error {
	// TODO Decode it into a symbol type.
	return InvalidTypeError{js.TypeSymbol, v.Type()}
}

// decodeArray decodes a JS array into the provided reflect.Value.
func decodeArray(x js.Value, v reflect.Value) error {
	jsLen := x.Length()

	switch v.Kind() {
	case reflect.Array:
		if jsLen != v.Len() {
			return InvalidArrayError{v.Len(), jsLen}
		}
	case reflect.Slice:
		newSlice := reflect.MakeSlice(v.Type(), jsLen, jsLen)
		v.Set(newSlice)
	default:
		return InvalidTypeError{js.TypeObject, v.Type()}
	}

	for i := 0; i < jsLen; i++ {
		err := FromJSValue(x.Index(i), v.Index(i).Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}

// decodeObject decodes a JS object into the provided reflect.Value.
func decodeObject(x js.Value, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		return decodeObjectIntoStruct(x, v)
	case reflect.Map:
		return decodeObjectIntoMap(x, v)
	default:
		return InvalidTypeError{js.TypeObject, v.Type()}
	}
}

// decodeObject decodes a JS object into the provided reflect.Value struct.
func decodeObjectIntoStruct(x js.Value, v reflect.Value) error {
	for i := 0; i < v.Type().NumField(); i++ {
		fieldType := v.Type().Field(i)
		if fieldType.PkgPath != "" {
			continue
		}

		name := fieldType.Name
		tagName, tagOK := fieldType.Tag.Lookup("wasm")

		if tagOK {
			if tagName == "-" {
				continue
			}
			name = tagName
		}

		err := decodeValue(x.Get(name), v.Field(i))
		if err != nil {
			if tagOK {
				return fmt.Errorf("in field %s (JS %s): %w", fieldType.Name, tagName, err)
			}
			return fmt.Errorf("in field %s: %w", fieldType.Name, err)
		}
	}

	return nil
}

func decodeObjectIntoMap(x js.Value, v reflect.Value) error {
	mapType := v.Type()
	keyType := mapType.Key()
	valType := mapType.Elem()

	switch keyType.Kind() {
	case reflect.String:
	case reflect.Interface:
		if keyType.NumMethod() != 0 {
			return InvalidTypeError{js.TypeObject, mapType}
		}
	default:
		return InvalidTypeError{js.TypeObject, mapType}
	}

	// TODO: Use Object API
	obj, err := Global().Get("Object")
	if err != nil {
		panic("Object not found")
	}

	var keys []string
	err = FromJSValue(obj.Call("keys", x), &keys)
	if err != nil {
		panic("Object.keys returned non-string-array.")
	}

	for _, k := range keys {
		valuePtr := reflect.New(valType).Interface()
		err := FromJSValue(x.Get(k), valuePtr)
		if err != nil {
			return err
		}

		v.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(valuePtr).Elem())
	}
	return nil
}

// decodeFunction decodes a JS function into the provided reflect.Value.
func decodeFunction(x js.Value, v reflect.Value) error {
	funcType := v.Type()
	outCount := funcType.NumOut()

	switch outCount {
	case 0, 1:
	case 2:
		if funcType.Out(1) != errorType {
			return ErrMultipleReturnValue
		}
	default:
		return ErrMultipleReturnValue
	}

	v.Set(reflect.MakeFunc(funcType, func(args []reflect.Value) []reflect.Value {
		argsJS := make([]interface{}, 0, len(args))
		for _, v := range args {
			argsJS = append(argsJS, ToJSValue(v.Interface()))
		}

		jsReturn := x.Invoke(argsJS...)
		if outCount == 0 {
			return []reflect.Value{}
		}

		returnPtr := reflect.New(funcType.Out(0)).Interface()
		err := FromJSValue(jsReturn, returnPtr)

		returnVal := reflect.ValueOf(returnPtr).Elem()
		if err != nil {
			if outCount == 1 {
				panic("error decoding JS return value: " + err.Error())
			}

			return []reflect.Value{returnVal, reflect.ValueOf(err)}
		}

		switch outCount {
		case 1:
			return []reflect.Value{returnVal}
		case 2:
			return []reflect.Value{returnVal, reflect.Zero(v.Type())}
		default:
			panic("unexpected amount of return values")
		}
	}))
	return nil
}

// createInterface creates a representation of the provided js.Value.
func createInterface(x js.Value) interface{} {
	switch x.Type() {
	case js.TypeUndefined, js.TypeNull:
		return nil
	case js.TypeBoolean:
		return x.Bool()
	case js.TypeNumber:
		return x.Float()
	case js.TypeString:
		return x.String()
	case js.TypeSymbol:
		// We can't convert it to a Go value in a meaningful way.
		return x
	case js.TypeObject:
		if isArray(x) {
			return createArray(x)
		}
		return createObject(x)
	case js.TypeFunction:
		var a func(...interface{}) (interface{}, error)
		err := FromJSValue(x, &a)
		if err != nil {
			panic("error creating function: " + err.Error())
		}
		return a
	default:
		panic("unknown JS type: " + x.Type().String())
	}
}

// createArray creates a slice of interface representing the js.Value.
func createArray(x js.Value) interface{} {
	result := make([]interface{}, x.Length())
	for i := range result {
		result[i] = createInterface(x.Index(i))
	}
	return result
}

// createObject creates a representation of the provided JS object.
func createObject(x js.Value) interface{} {
	// TODO: Use Object API
	obj, err := Global().Get("Object")
	if err != nil {
		panic("Object not found")
	}

	var keys []string
	err = FromJSValue(obj.Call("keys", x), &keys)
	if err != nil {
		panic("Object.keys returned non-string-array.")
	}

	result := make(map[string]interface{}, len(keys))
	for _, v := range keys {
		result[v] = createInterface(x.Get(v))
	}
	return result
}

// isArray calls the JS function Array.isArray to check if the provided js.Value is an array.
func isArray(x js.Value) bool {
	arr, err := Global().Get("Array")
	if err != nil {
		panic("Array not found")
	}

	return arr.Call("isArray", x).Bool()
}

// initializePointerIfNil checks if the pointer is nil and initializes it as necessary.
func initializePointerIfNil(v reflect.Value) {
	if v.Kind() != reflect.Ptr {
		return
	}
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	initializePointerIfNil(v.Elem())
}
