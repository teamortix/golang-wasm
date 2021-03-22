import wasm from './api/main.go';

const { hello, helloName } = wasm;

const value = document.getElementById("value");
const input = document.getElementById("input");
const funcValue = document.getElementById("funcValue");

const run = async () => {
    value.innerText = await hello();

    funcValue.innerText = await helloName(input.value);
    input.addEventListener("keyup", async (e) => {
        funcValue.innerText = await helloName(e.target.value);
    })
}

run()
