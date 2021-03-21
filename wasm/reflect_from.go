package wasm

import (
	"fmt"
	"syscall/js"
)

// FromJSValue converts a given js.Value to the Go equivalent.
// The new value of 'out' is undefined if FromJSValue returns an error.
func FromJSValue(x js.Value, out interface{}) error {
	// TODO
	return fmt.Errorf("unimplemented")
}
