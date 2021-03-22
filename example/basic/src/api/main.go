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

func setup() {
	fmt.Println("golang-wasm initialized")

	js.Global()
}
