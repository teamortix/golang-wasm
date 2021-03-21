package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Hello from go-mod-wasm!")
	js.Global().Get("__go_wasm__").Set("__ready__", true)

	c := make(chan bool, 0) // in Go Wasm, the program may not exit
	<-c
}
