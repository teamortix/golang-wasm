const g = global || window || self
if (!g.__go_wasm__) {
    g.__go_wasm__ = {};
}

const maxTime = new Date()
maxTime.setSeconds(maxTime.getSeconds() + 3) // if js does not initialize after 3 seconds, we allow it to start anyhow and print a warning

const bridge = g.__go_wasm__;

function sleep() {
    return new Promise(requestAnimationFrame)
}
export default function (getBytes) {

    async function init() {
        const go = new g.Go();
        let bytes = await getBytes
        let result = await WebAssembly.instantiate(bytes, go.importObject);
        go.run(result.instance);
        setTimeout(() => {
            if (bridge.__ready__ !== true) {
                console.warn("Golang Wasm Bridge (__go_wasm__.__ready__) still not true after max time")
            }
        }, 3 * 1000)
    }

    init();


    let proxy = new Proxy(
        {},
        {
            get: (_, key) => {
                return (...args) => {
                    return new Promise(async (res, rej) => {
                        while (bridge.__ready__ !== true) {
                            await sleep()
                        }

                        if (typeof bridge[key] !== 'function') {
                            res(bridge[key]);
                            return;
                        }

                        const returnObj = bridge[key].apply(undefined, args)
                        if (returnObj.error) {
                            rej(returnObj.error)
                        } else {
                            res(returnObj.result)
                        }
                    })
                };
            }
        }
    );

    return proxy;
}
