package wasm

import "syscall/js"

// Promise is an instance of a JS promise.
// The zero value of this struct is not a valid Promise.
type Promise struct {
	Object
}

// FromJSValue turns a JS value to a Promise.
func (p *Promise) FromJSValue(value js.Value) error {
	var err error
	p.Object, err = NewObject(value)
	return err
}

// NewPromise returns a promise that is fulfilled or rejected when the provided handler returns.
// The handler is spawned in its own goroutine.
func NewPromise(handler func() (interface{}, error)) Promise {
	resultChan := make(chan interface{})
	errChan := make(chan error)

	// Invoke the handler in a new goroutine.
	go func() {
		result, err := handler()
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	// Create a JS promise handler.
	var jsHandler js.Func
	jsHandler = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			panic("not enough arguments are passed to the Promise constructor handler")
		}

		resolve := args[0]
		reject := args[1]

		if resolve.Type() != js.TypeFunction || reject.Type() != js.TypeFunction {
			panic("invalid type passed to Promise constructor handler")
		}

		go func() {
			select {
			case r := <-resultChan:
				resolve.Invoke(ToJSValue(r))
			case err := <-errChan:
				reject.Invoke(NewError(err))
			}

			// Free up resources now that we are done.
			jsHandler.Release()
		}()

		return nil
	})

	promise, err := Global().Expect(js.TypeFunction, "Promise")
	if err != nil {
		panic("Promise constructor not found")
	}

	return mustJSValueToPromise(promise.New(jsHandler))
}

// PromiseAll creates a promise that is fulfilled when all the provided promises have been fulfilled.
// The promise is rejected when any of the promises provided rejects.
// It is implemented by calling Promise.all on JS.
func PromiseAll(promise ...Promise) Promise {
	promiseAll, err := Global().Expect(js.TypeFunction, "Promise", "all")
	if err != nil {
		panic("Promise.all not found")
	}

	pInterface := make([]interface{}, 0, len(promise))
	for _, v := range promise {
		pInterface = append(pInterface, v)
	}

	return mustJSValueToPromise(promiseAll.Invoke(pInterface))
}

// PromiseAllSettled creates a promise that is fulfilled when all the provided promises have been fulfilled or rejected.
// It is implemented by calling Promise.allSettled on JS.
func PromiseAllSettled(promise ...Promise) Promise {
	promiseAllSettled, err := Global().Expect(js.TypeFunction, "Promise", "allSettled")
	if err != nil {
		panic("Promise.allSettled not found")
	}

	pInterface := make([]interface{}, 0, len(promise))
	for _, v := range promise {
		pInterface = append(pInterface, v)
	}

	return mustJSValueToPromise(promiseAllSettled.Invoke(pInterface))
}

// PromiseAny creates a promise that is fulfilled when any of the provided promises have been fulfilled.
// The promise is rejected when all of the provided promises gets rejected.
// It is implemented by calling Promise.any on JS.
func PromiseAny(promise ...Promise) Promise {
	promiseAny, err := Global().Expect(js.TypeFunction, "Promise", "any")
	if err != nil {
		panic("Promise.any not found")
	}

	pInterface := make([]interface{}, 0, len(promise))
	for _, v := range promise {
		pInterface = append(pInterface, v)
	}

	return mustJSValueToPromise(promiseAny.Invoke(pInterface))
}

// PromiseRace creates a promise that is fulfilled or rejected when one of the provided promises fulfill or reject.
// It is implemented by calling Promise.race on JS.
func PromiseRace(promise ...Promise) Promise {
	promiseRace, err := Global().Expect(js.TypeFunction, "Promise", "race")
	if err != nil {
		panic("Promise.race not found")
	}

	pInterface := make([]interface{}, 0, len(promise))
	for _, v := range promise {
		pInterface = append(pInterface, v)
	}

	return mustJSValueToPromise(promiseRace.Invoke(pInterface))
}

func mustJSValueToPromise(v js.Value) Promise {
	var p Promise
	err := p.FromJSValue(v)
	if err != nil {
		panic("Expected a Promise from JS standard library")
	}

	return p
}
