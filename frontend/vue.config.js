// vue.config.js
const webpack = require('webpack');
const isProd = process.env.NODE_ENV === "production";

module.exports = {
    outputDir: "../backend/static",
    configureWebpack: {
        devServer: {
            proxy: {
                "/api": {
                    "target": "http://localhost:8000",
                    "secure": false,
                    "changeOrigin": true
                },
                "/svg": {
                    "target": "http://localhost:8000",
                    "secure": false,
                    "changeOrigin": true
                },
                "/js/popper.min.js": {
                    "target": "http://localhost:8000",
                    "secure": false,
                    "changeOrigin": true
                },
                "/js/bootstrap.min.js": {
                    "target": "http://localhost:8000",
                    "secure": false,
                    "changeOrigin": true
                },
                "/css/bootstrap.min.css": {
                    "target": "http://localhost:8000",
                    "secure": false,
                    "changeOrigin": true
                },
            }
        }
    }
};

