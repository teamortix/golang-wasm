package wasm_test

import (
	"syscall/js"
	"testing"

	"github.com/teamortix/golang-wasm/wasm"
)

// Magic values to communicate with the JS library.
const (
	globalIdent     = "__go_wasm__"
	proxyName       = "__proxy__"
	readyHint       = "__ready__"
	funcWrapperName = "__wrapper__"
)

var (
	bridge      wasm.Object
	funcWrapper js.Value
	proxy       js.Value
)

func TestMain(t *testing.M) {
	bridgeJS, err := wasm.Global().Get(globalIdent)
	if err != nil {
		panic("JS wrapper " + globalIdent + " not found")
	}

	bridge, err = wasm.NewObject(bridgeJS)
	if err != nil {
		panic("JS wrapper " + globalIdent + " is not an object")
	}

	funcWrapper, err = bridge.Get(funcWrapperName)
	if err != nil {
		panic("JS wrapper " + globalIdent + "." + funcWrapperName + " not found")
	}

	proxy, err = bridge.Get(proxyName)
	if err != nil {
		panic("JS proxy " + globalIdent + "." + proxyName + " not found")
	}

	t.Run()
}

func TestSetupCorrectly(t *testing.T) {
	if funcWrapper.Type() != js.TypeFunction {
		t.Errorf("expected wrapper to return %s, instead returned %s", js.TypeFunction, funcWrapper.Type())
		return
	}

	if proxy.Type() != js.TypeObject {
		t.Errorf("expected proxy to return %s, instead returned %s", js.TypeObject, proxy.Type())
		return
	}
}

func TestReady(t *testing.T) {
	ready, err := bridge.Expect(js.TypeBoolean, readyHint)
	if ready.Truthy() {
		t.Errorf("expected ready value to be falsy before call to wasm.Ready()")
		return
	}

	wasm.Ready()
	ready, err = bridge.Expect(js.TypeBoolean, readyHint)
	if err != nil {
		t.Errorf("expected ready value in bridge to be boolean")
		return
	}
	if !ready.Bool() {
		t.Errorf("expected ready value to be true after call to wasm.Ready()")
		return
	}
}
