const fs = require("fs")
const { execFileSync } = require("child_process")
const path = require("path")
const { lookpath } = require("lookpath")

const exists = async (dir, file) => {
    return new Promise((res, rej) => {
        fs.access(path.join(dir, file), fs.constants.F_OK, (err) => {
            if (err) {
                return res(false)
            }
            return res(true)
        })
    })
}

module.exports = function (source) {
    const cb = this.async()

    const goBin = lookpath("go");
    if (!goBin) {
        return cb(new Error("go bin not found in path."));
    }

    if (!process.env.GOROOT) {
        return cb(new Error("Could not find GOROOT in environment.\n" +
            "Please try adding this to your script:\nGOROOT=`go env GOROOT` npm run ..."))
    }

    const parent = path.dirname(this.resourcePath)
    const outFile = this.resourcePath.slice(0, -2) + "wasm"
    let modDir = parent
    let found = false;

    const opts = {
        cwd: parent,
        env: {
            GOPATH: process.env.GOPATH,
            GOROOT: process.env.GOROOT,
            GOCACHE: path.join(__dirname, ".gocache"),
            GOOS: "js",
            GOARCH: "wasm",
        }
    };

    (async () => {
        const root = path.resolve(path.sep)
        while (path.resolve(modDir) != root) {
            if (!(await exists(modDir, 'go.mod'))) {
                modDir = path.join(modDir, '..');
            } else {
                found = true;
                break;
            }
        }
        if (!found) {
            return cb(new Error("Could not find go.mod in any parent directory of " + this.resourcePath))
        }

        try {
            execFileSync("go", ["build", "-o", outFile, parent], opts)
        } catch (e) {
            return cb(e)
        }

        const wasmOrigPath = path.join(process.env.GOROOT, "misc", "wasm", "wasm_exec.js")
        const wasmEmitPath = path.join(__dirname, 'wasm_exec.js')
        if (!(await exists(__dirname, 'wasm_exec.js'))) {
            fs.copyFileSync(wasmOrigPath, wasmEmitPath)
        }
        const contents = fs.readFileSync(outFile)
        fs.unlinkSync(outFile)

        const emitPath = path.basename(outFile)
        this.emitFile(emitPath, contents)
        this.addContextDependency(modDir)

        cb(null,
            `require('!${wasmEmitPath}')
import goWasm from '${path.join(__dirname, 'bridge.js')}';

const wasm = fetch('${emitPath}').then(response => response.arrayBuffer())
export default goWasm(wasm)`)
    })()
}