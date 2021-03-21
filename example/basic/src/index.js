import wasm from './api/main.go';

const { hello, helloName } = wasm;

(async () => {
    console.log(await hello());
    console.log(await helloName("world"));
})()