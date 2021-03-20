package main

import (
	"fmt"

	"gitea.teamortix.com/Team-Ortix/go-mod-wasm/wasm"
)

func main() {
	c := make(chan bool, 0)
	fmt.Println("Hello Go!")
	wasm.Ready()
	<-c
}
