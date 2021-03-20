package wasm

import (
	"fmt"
	"syscall/js"
)

// TypeMismatchError is returned when a function is called with a js.Value that has the incorrect type.
type TypeMismatchError struct {
	Expected js.Type
	Actual   js.Type
}

func (e TypeMismatchError) Error() string {
	return fmt.Sprintf("expected %v type, got %v type instead", e.Expected, e.Actual)
}

// Global returns the global object as a Object.
// If the global object is not an object, it panics.
func Global() Object {
	global, err := NewObject(js.Global())
	if err != nil {
		panic(err)
	}
	return global
}

// Object is a statically typed Object instance of js.Value.
// It should be instantiated with NewObject where it is checked for type instead of directly.
// Calling methods on a zero Object is undefined behaviour.
type Object struct {
	value js.Value
}

// NewObject instantiates a new Object with the provided js.Value.
// If the js.Value is not an Object, it returns a TypeMismatchError.
func NewObject(raw js.Value) (Object, error) {
	if raw.Type() != js.TypeObject {
		return Object{}, TypeMismatchError{
			Expected: js.TypeObject,
			Actual:   raw.Type(),
		}
	}

	return Object{raw}, nil
}

// Get recursively gets the Object's properties, returning a TypeMismatchError if it encounters a non-object while
// descending through the object.
func (o Object) Get(path ...string) (js.Value, error) {
	current := o.value
	for _, v := range path {
		if current.Type() != js.TypeObject {
			return js.Value{}, TypeMismatchError{
				Expected: js.TypeObject,
				Actual:   current.Type(),
			}
		}

		current = current.Get(v)
	}
	return current, nil
}

// Expect is a helper function that calls Get and checks the type of the final result.
// It returns a TypeMismatchError if a non-object is encountered while descending the path or the final type does not
// match with the provided expected type.
func (o Object) Expect(expectedType js.Type, path ...string) (js.Value, error) {
	value, err := o.Get(path...)
	if err != nil {
		return js.Value{}, err
	}

	if value.Type() != expectedType {
		return js.Value{}, TypeMismatchError{
			Expected: expectedType,
			Actual:   value.Type(),
		}
	}

	return value, nil
}

// Delete removes property p from the object.
func (o Object) Delete(p string) {
	o.value.Delete(p)
}

// Equal checks if the object is equal to another value.
// It is equivalent to JS's === operator.
func (o Object) Equal(v js.Value) bool {
	return o.value.Equal(v)
}

// Index indexes into the object.
func (o Object) Index(i int) js.Value {
	return o.value.Index(i)
}

// InstanceOf implements the instanceof operator in JavaScript.
// If t is not a constructor, this function returns false.
func (o Object) InstanceOf(t js.Value) bool {
	if t.Type() != js.TypeFunction {
		return false
	}
	return o.value.InstanceOf(t)
}

// JSValue implements the js.Wrapper interface.
func (o Object) JSValue() js.Value {
	return o.value
}

// Length returns the "length" property of the object.
func (o Object) Length() int {
	return o.value.Length()
}

// Set sets the property p to the value of ToJSValue(x).
func (o Object) Set(p string, x interface{}) {
	o.value.Set(p, ToJSValue(x))
}

// SetIndex sets the index i to the value of ToJSValue(x).
func (o Object) SetIndex(i int, x interface{}) {
	o.value.SetIndex(i, ToJSValue(x))
}

// String returns the object marshalled as a JSON string for debugging purposes.
func (o Object) String() string {
	stringify, err := Global().Expect(js.TypeFunction, "JSON", "stringify")
	if err != nil {
		panic(err)
	}

	jsonStr := stringify.Invoke(o)
	if jsonStr.Type() != js.TypeString {
		panic("JSON.stringify returned a " + jsonStr.Type().String())
	}

	return jsonStr.String()
}
