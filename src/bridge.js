const g = global || window || self;
if (!g.__go_wasm__) {
    g.__go_wasm__ = {};
}

const maxTime = 3 * 1000;

const bridge = g.__go_wasm__;

/**
 * Wrapper is used by Go to run all Go functions in JS.
 * Go functions always return an object of the following spec:
 * {
 *  result:  undefined | any         // undefined when error is returned, or function returns undefined
 *  error:       Error | undefined   // undefined when no error is present
 * }
 */
function wrapper(goFunc) {
    return (...args) => {
        const result = goFunc.apply(undefined, args);
        if (result.error instanceof Error) {
            throw result.error;
        }
        return result.result;
    }
}

function sleep() {
    return new Promise((res) => {
        requestAnimationFrame(() => res())
        setTimeout(() => {
            res()
        }, 50);
    });
}
export default function (getBytes) {
    let proxy;

    async function init() {
        bridge.__wrapper__ = wrapper

        const go = new g.Go();
        let bytes = await getBytes;
        let result = await WebAssembly.instantiate(bytes, go.importObject);
        go.run(result.instance);
    }

    init();
    setTimeout(() => {
        if (bridge.__ready__ !== true) {
            console.warn("Golang Wasm Bridge (__go_wasm__.__ready__) still not true after max time");
        }
    }, maxTime);


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

                        try {
                            res(bridge[key].apply(undefined, args));
                        } catch (e) {
                            rej(e)
                        }
                    })
                };
            }
        }
    );

    bridge.__proxy__ = proxy
    return proxy;
}
