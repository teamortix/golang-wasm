# Golang-WASM Architecture

> Please note that Golang-WASM is still in its initial stages of development and is not at a stable version yet. SemVer will be used when v1.0.0 is released and changes are expected to be potentially breaking and unstable.

This file gives a mile-high view of the Golang-WASM project and how the different parts interact with each other. We strongly recommend all contributors to read this, as well as [CONTRIBUTING.md](./CONTRIBUTING.md) before getting started with development. 

Go-WASM has 2 primary parts to the project.

* The Go WASM [bindings](./wasm) and `syscall/js` superset to improve type safety when working in Go. 
* The JavaScript caller and wrappers, and the Webpack loader a which can be found in the [src](./src) folder.

## Go WASM Bindings

The Go WASM bindings has three main goals,
1. Make type checking a non-issue for development.
2. Seamlessly convert the programming style from Go to JS.
3. (WIP) Implement a type safe version of the entire DOM API.

Solving the first two require a healthy usage of reflection and their implementation can be found mostly in [reflect_to.go](./wasm/reflect_to.go) and [reflect_from.go](./wasm/reflect_from.go). 

Interfacing from functions that return errors to Promise resolves and rejections is mostly handled within [function.go](./wasm/function.go). It is worth noting how error handling works when working with calling functions from Go in a type safe manner in [reflect_from](./wasm/reflect_from.go) within `decodeFunc`.

### DOM API

This is still a work in progress. However, basic parts of it have already been implemented. The implementation for the [Promise API](./wasm/promise.go) demonstrates what we have in mind for the rest of the API. The goal of this project is not to dump everything 1:1. If you want to use something like that, you can use a computer generated version [here](https://github.com/brettlangdon/go-dom).



## JavaScript + Webpack

The main goal of the JavaScript hook for the project is to seamlessly link Go and JavaScript together and overcome some of the limitations of native Go WASM through wrappers. Read the implementation [here](./src/bridge.js)

The main difference between writing JavaScript code and Go code is the interactivity with asynchronous code, and error handling. To solve this issue, calling all Go code from JS goes through a [Proxy](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Proxy). Furthermore, all Go functions that are called from JavaScript are returned via a wrapper functions. Read more about this in the documentation.

### Webpack Loader

The Webpack loader is relatively simple to use. To use it, it is only important for developers to add `GOROOT` to their environment variables when running the command.

The Webpack loader searches and finds a `go.mod` file in any of the parent directories of an imported Go file. The loader requires the presence of such a file. Webpack then watches for all file changes inside the go.mod directory and **compiles your entire project**, as long as one of the files are imported in the Go project.  Read the implementation [here](./src/index.js)

The Webpack Loader for Golang-Wasm is not in its ideal state. Currently, It does not seamlessly work with Webpack supersets, such as [Craco](https://github.com/gsoft-inc/craco/issues/268). 