import wasm from './main.go';

(async () => {
    console.log(wasm)
    console.log(await wasm.test(), "..")
})()