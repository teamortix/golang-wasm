const fs = require("fs/promises");
const util = require("util");
const execFile = util.promisify(require("child_process").execFile);
const path = require("path");
const { lookpath } = require("lookpath");

module.exports = function (source) {
    const cb = this.async();

    const goBin = lookpath("go");
    if (!goBin) {
        return cb(new Error("go bin not found in path."));
    }

    if (!process.env.GOROOT) {
        return cb(new Error("Could not find GOROOT in environment.\n" +
            "Please try adding this to your script:\n" +
            "GOROOT=`go env GOROOT` npm run ..."));
    }

    const parent = path.dirname(this.resourcePath);
    const outFile = this.resourcePath.slice(0, -2) + "wasm";
    let modDir = parent;

    const opts = {
        cwd: parent,
        env: {
            GOPATH: process.env.GOPATH || path.join(process.env.HOME, "go"),
            GOROOT: process.env.GOROOT,
            GOCACHE: path.join(__dirname, ".gocache"),
            GOOS: "js",
            GOARCH: "wasm",
        },
    };

    (async () => {
        let found = false;
        const root = path.resolve(path.sep);
        while (path.resolve(modDir) != root) {
            found = await fs.access(path.join(modDir, 'go.mod')).then(() => true).catch(() => false);
            if (found) {
                break;
            }
            modDir = path.join(modDir, "..");
        }

        if (!found) {
            return cb(new Error("Could not find go.mod in any parent directory of " + this.resourcePath));
        }

        const wasmOrigPath = path.join(process.env.GOROOT, "misc", "wasm", "wasm_exec.js");
        const wasmSavePath = path.join(__dirname, 'wasm_exec.js');
        const errorPaths = ["\t" + wasmOrigPath, "\t" + wasmSavePath];
        if (!(await fs.access(wasmOrigPath).then(() => true).catch(() => false)) &&
            !(await fs.access(wasmSavePath).then(() => true).catch(() => false))) {
            return cb(new Error("Could not find wasm_exec.js file. Invalid GOROOT? Searched paths:\n" +
                errorPaths.join(",\n") + "\n"));
        }


        const res = await execFile("go", ["build", "-o", outFile, parent], opts)
            .then(() => true)
            .catch(e => e);
        if (res instanceof Error) {
            return cb(e);
        }

        found = await fs.access(wasmSavePath).then(() => true).catch(() => false);
        if (!found) fs.copyFile(wasmOrigPath, wasmSavePath);

        const contents = await fs.readFile(outFile);
        fs.unlink(outFile);

        const emitPath = path.basename(outFile);
        this.emitFile(emitPath, contents);
        this.addContextDependency(modDir);

        cb(null,
            `require('!${wasmSavePath}');
import goWasm from '${path.join(__dirname, 'bridge.js')}';

const wasm = fetch('${emitPath}').then(response => response.arrayBuffer());
export default goWasm(wasm);`);
    })();
}
