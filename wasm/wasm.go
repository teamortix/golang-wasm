package wasm

import "syscall/js"

// Magic values to communicate with the JS library.
const (
	globalIdent     = "__go_wasm__"
	readyHint       = "__ready__"
	funcWrapperName = "__wrapper__"
)

var (
	bridge      Object
	funcWrapper js.Value
)

func init() {
	bridgeJS, err := Global().Get(globalIdent)
	if err != nil {
		panic("JS wrapper " + globalIdent + " not found")
	}

	bridge, err = NewObject(bridgeJS)
	if err != nil {
		panic("JS wrapper " + globalIdent + " is not an object")
	}

	funcWrapper, err = bridge.Get(funcWrapperName)
	if err != nil {
		panic("JS wrapper " + globalIdent + "." + funcWrapperName + " not found")
	}
}

// Ready notifies the JS bridge that the WASM is ready.
// It should be called when every value and function is exposed.
func Ready() {
	Expose(readyHint, true)
}

// Expose exposes a copy of the provided value in JS.
func Expose(property string, x interface{}) {
	bridge.Set(property, x)
}
