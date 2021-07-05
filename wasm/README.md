<p align="center">
    <a href="https://github.com/teamortix/golang-wasm">
        <img src="../banner.png">
    </a>
</p>
<p align="center">A bridge and bindings for JS DOM API with Go WebAssembly.</p>
<p align="center">Written by Team Ortix - <a href="https://github.com/hhhapz">Hamza Ali</a> and <a href="https://github.com/chanbakjsd/">Chan Wen Xu</a>.</p>
<p align="center">
    <a href="https://pkg.go.dev/github.com/teamortix/golang-wasm/wasm">
        <img src="https://pkg.go.dev/badge/github.com/teamortix/golang-wasm/wasm.svg" alt="Go Reference">
    </a>
    <a href="https://goreportcard.com/report/github.com/teamortix/golang-wasm">
        <img src="https://goreportcard.com/badge/github.com/teamortix/golang-wasm" alt="Go Report Card">
    </a>
    <br>
    <br>
</p>

## [Documentation Available Here](../README.md)

```
GOOS=js GOARCH=wasm go get -u github.com/teamortix/golang-wasm/wasm
```
> To run tests, run [test.sh](./testing/test.sh)

```bash
sh testing/test.sh . # Test all files in directory.

sh testing/wasm_test.go # Test only a single file.
```

---
Created by [hhhapz](https://github.com/hhhapz) and [chanbakjsd](https://github.com/chanbakjsd)
