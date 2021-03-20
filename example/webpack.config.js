const path = require('path');

module.exports = {
    entry: './src/index.js',
    mode: "production",
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'dist'),
    },
    resolve: {
        extensions: [".go", ".tsx", ".ts", ".js"],
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
            "crypto-browserify": require.resolve('crypto-browserify'), //if you want to use this module also don't forget npm i crypto-browserify 
        }
    },
    module: {
        rules: [
            {
                test: /\.go$/,
                use: [
                    {
                        loader: path.resolve(__dirname, '../index.js')
                    }
                ]
            }
        ]
    },
};