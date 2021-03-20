package wasm

import "syscall/js"

// NewError returns a JS Error with the provided Go error's error message.
func NewError(goErr error) js.Value {
	errConstructor, err := Global().Expect(js.TypeFunction, "Error")
	if err != nil {
		panic("Error constructor not found")
	}

	return errConstructor.New(goErr.Error())
}
