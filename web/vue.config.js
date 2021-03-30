'use strict';
const path = require('path');
const FilemanagerWebpackPlugin = require('filemanager-webpack-plugin');
const resolve = function(dir) {
    return path.join(__dirname, dir);
};

module.exports = {
    publicPath: '/page',
    devServer: {
        port: 8085,
        open: true,
        hot: true,
        proxy: {
            '^/api': {
                target: `http://127.0.0.1:8086/`, //代理到 目标路径
                changeOrigin: true,
            },
            '^/portal': {
                target: `http://127.0.0.1:8086/`, //代理到 目标路径
                changeOrigin: true,
            },
        },
        disableHostCheck: true,
        historyApiFallback: {
            rewrites: [
                {
                    from: /\/page\/dashboard/g,
                    to: '/page/dashboard.html',
                },
                {
                    from: /\/page\/portal/g,
                    to: '/page/portal.html',
                },
                {
                    from: /\/page\/proxy/g,
                    to: '/page/proxy.html',
                },
            ],
        },
    },
    configureWebpack: {
        plugins: [
            new FilemanagerWebpackPlugin({
                events: {
                    onStart: {
                        delete: [
                            {
                                source: resolve('../cmd/web/dist'),
                                options: {
                                    force: true,
                                },
                            },
                        ],
                    },
                    onEnd: {
                        copy: [
                            {
                                source: resolve('dist'),
                                destination: resolve('../cmd/web/dist'),
                            },
                        ],
                    },
                },
            }),
        ],
    },
    chainWebpack: (config) => {
        config.resolve.alias
            .set('@', resolve('src'))
            .set('@dashboard', resolve('src/pages/dashboard'))
            .set('@portal', resolve('src/pages/portal'))
            .set('@proxy', resolve('src/pages/proxy'));
    },
    pages: {
        dashboard: {
            entry: 'src/pages/dashboard/main.js',
            template: 'public/dashboard.html',
            filename: 'dashboard.html',
            title: 'User Growth管理台',
        },
        portal: {
            entry: 'src/pages/portal/main.js',
            template: 'public/portal.html',
            filename: 'portal.html',
            title: 'User Growth',
        },
        proxy: {
            entry: 'src/pages/proxy/main.js',
            template: 'public/proxy.html',
            filename: 'proxy.html',
            title: 'User Growth Proxy',
        },
    },
};
