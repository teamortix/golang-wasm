const path = require('path');

module.exports = {
    entry: './src/index.js',
    devServer: {
		static: {
        	directory: path.resolve(__dirname, 'dist'),
		},
        compress: true,
        port: 3000,
    },
    mode: "development",
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'dist'),
    },
    resolve: {
        extensions: [".js", ".go"],
        fallback: {
            "fs": false,
            "os": false,
            "util": false,
            "tls": false,
            "net": false,
            "path": false,
            "zlib": false,
            "http": false,
            "https": false,
            "stream": false,
            "crypto": false,
        }
    },
    module: {
        rules: [
            {
                test: /\.go$/,
                use: [
                    {
                        loader: path.resolve(__dirname, '../../src/index.js')
                    }
                ]
            }
        ]
    },
    performance: {
        assetFilter: (file) => {
            return !/(\.wasm|.map)$/.test(file)
        }
    },
    ignoreWarnings: [
        {
            module: /wasm_exec.js$/
        }
    ]
};
