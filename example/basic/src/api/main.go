package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Hello from go-mod-wasm!")
	setup()

	c := make(chan bool, 0) // To use anything from Go WASM, the program may not exit.
	<-c
}

const hello = "Sample value"

func helloName(_ js.Value, args []js.Value) interface{} {
	return fmt.Sprintf("Hello, %s!", args[0].String())
}

func setup() {
	bridge := js.Global().Get("__go_wasm__")

	bridge.Set("__ready__", true)

	bridge.Set("hello", hello)
	bridge.Set("helloName", js.FuncOf(helloName))

	js.Global()
}
