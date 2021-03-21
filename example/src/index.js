import wasm from './api/main.go';

(async () => {
    console.log(await wasm.__ready__())
})()