const g = global || window || self;
if (!g.__go_wasm__) {
    g.__go_wasm__ = {};
}

const maxTime = 3 * 1000;

const bridge = g.__go_wasm__;

function sleep() {
    return new Promise(requestAnimationFrame);
}
export default function (getBytes) {
    let proxy;

    async function init() {
        const go = new g.Go();
        let bytes = await getBytes;
        let result = await WebAssembly.instantiate(bytes, go.importObject);
        go.run(result.instance);
        bridge.__proxy__ = proxy
        setTimeout(() => {
            if (bridge.__ready__ !== true) {
                console.warn("Golang Wasm Bridge (__go_wasm__.__ready__) still not true after max time");
            }
        }, maxTime);
    }

    init();


    proxy = new Proxy(
        {},
        {
            get: (_, key) => {
                return (...args) => {
                    return new Promise(async (res, rej) => {
                        while (bridge.__ready__ !== true) {
                            await sleep();
                        }

                        if (typeof bridge[key] !== 'function') {
                            res(bridge[key]);
                            return;
                        }

                        const returnObj = bridge[key].apply(undefined, args);
                        if (returnObj.error instanceof Error) {
                            return rej(returnObj.error)
                        }

                        if (returnObj.result) return res(returnObj.result);

                        return res(returnObj)
                    })
                };
            }
        }
    );

    return proxy;
}
