import wasm from './api/main.go';

const { hello, helloName, divide } = wasm;

const value = document.getElementById("hello");

const inputName = document.getElementById("inputName");
const name = document.getElementById("name");

const inputX = document.getElementById("inputX");
const inputY = document.getElementById("inputY");
const divideElem = document.getElementById("divide");


const doDivision = async (x, y) => divideElem.innerText =
    await divide(Number(x), Number(y))
        .catch(err => err.toString());

const run = async () => {
    value.innerText = await hello();

    name.innerText = await helloName(inputName.value);
    doDivision(inputX.value, inputY.value);

    inputName.addEventListener("keyup", async (e) => name.innerText = await helloName(e.target.value));
    inputX.addEventListener("change", (e) => doDivision(e.target.value, inputY.value));
    inputY.addEventListener("change", (e) => doDivision(inputX.value, e.target.value, inputY.value));
}

run()
