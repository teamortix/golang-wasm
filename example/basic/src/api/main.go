package main

import (
	"fmt"
	"syscall/js"
)

const hello = "Hello!"

// helloName's first value is JavaScript's `this`.
// However, the way that the JS bridge is written, it will always be JavaScript's undefined.
//
// If returning a non-nil error value, the resulting promise will be rejected by API consumers.
// The rejected value will JavaScript's Error, with the message being the go error's message.
//
// See other examples which use the Go WASM bridge api, which show more flexibility and type safety when interacting
// with JavaScript.
func helloName(_ js.Value, args []js.Value) (interface{}, error) {
	return fmt.Sprintf("Hello, %s!", args[0].String()), nil
}

func main() {
	fmt.Println("go-mod-wasm initialized")

	setFunc("helloName", helloName)
	setValue("hello", hello)
	ready()
}
