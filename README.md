<p align="center">
    <a href="https://github.com/teamortix/golang-wasm">
        <img src="./banner.png">
    </a>
</p>
<p align="center">A bridge and bindings for JS DOM API with Go WebAssembly.</p>
<p align="center">Written by Team Ortix - <a href="https://github.com/hhhapz">Hamza Ali</a> and <a href="https://github.com/chanbakjsd/">Chan Wen Xu</a></p>
<p align="center">
    <img src="https://godoc.org/github.com/teamortix/golang-wasm/wasm?status.svg">
    <img src="https://goreportcard.com/badge/github.com/teamortix/golang-wasm/wasm">
    <br>
    <br>
</p>

```
GOOS=js GOARCH=wasm go get -u github.com/teamortix/golang-wasm/wasm
```
```bash
npm install golang-wasm
```

# ⚠️ The documentation is still work in progress.

### Go API Documentation

[Reference site](https://pkg.go.dev/github.com/teamortix/golang-wasm)


## Why Golang-WASM?

Golang-WASM provides a simple idiomatic, and comprehensive (soon™️) API and bindings for working with WebAssembly.

Golang-WASM also comes with a [webpack](https://npmjs.com/golang-wasm) loader that wraps the entire API so that it is idiomatic for JavaScript develoeprs as well as Go developers.

Here is a small snippet:
```go
// Automatically handled with Promise rejects when returning an error!
func divide(x int, y int) (int, error) {
    if y == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return x / y, nil
}

func main() {
    wasm.Expose("divide", divide)
    wasm.Ready()
}
```
```js
import { divide } from "main.go"

const result = await divide(6, 2)
console.log(result) // 3

// Unhandled rejection in promise: cannot divide by zero
const error = await divide(6, 0) 
```

When using the webpack loader, everything is bundled for you, and you can directly import the Go file.

> Note: the webpack loader expects you to have a valid Go installation on your system, and a valid GOROOT passed.

## JS Interop

### Examples



### Working with functions

### Auto type casting

### Working with errors

### DOM API

### How it works

---

## Configuration

### Webpack Configuration

### Hotcode Reload

## FAQ

### Is it possible to use multiple instances of Web Assembly in the same project

At the moment, this is not supported.

### License

MIT

---
Created by [hhhapz](https://github.com/hhhapz) and [chanbakjsd](https://github.com/chanbakjsd)