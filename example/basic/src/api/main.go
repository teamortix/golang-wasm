package main

import (
	"errors"
	"fmt"

	"github.com/teamortix/golang-wasm/wasm"
)

const hello = "Hello!"

// The Golang-WASM API automatically ports calls to this function calls from JavaScript over to call this function.
// This helps gain a level of type safety that is otherwise not possible when using syscall/js.
// If this function is called from JavaScript with invalid arguments, the promise is simply rejected.
// Explore further how error handling works in the documentation.
func helloName(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

// This call is automatically converted to a promise resolve or rejection.
// This is similar to the example in the README, however it uses float64 instead so division results are not truncated.
//
// If this function's second return value is a non-value nil, Golang-WASM will reject the promise of the JS call using
// the error message from calling Error().
func divide(x float64, y float64) (float64, error) {
	if y == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return x / y, nil
}

func main() {
	fmt.Println("golang-wasm initialized")

	wasm.Expose("hello", hello)
	wasm.Expose("helloName", helloName)
	wasm.Expose("divide", divide)
	wasm.Ready()
	<-make(chan bool) // To use anything from Go WASM, the program may not exit.
}
