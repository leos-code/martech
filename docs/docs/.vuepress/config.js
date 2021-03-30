'use strict';
const path = require('path');
const FilemanagerWebpackPlugin = require('filemanager-webpack-plugin');
const fse = require('fs-extra');
const resolve = function(dir) {
    return path.join(__dirname, dir);
};
module.exports = {
    base: process.env.VUEPRESS_BASE || '/docs/',
    title: 'User Growth',
    head: [['link', { ref: 'icon', href: '/logo.png' }]],
    plugins: ['vuepress-plugin-mermaidjs'],
    dest: resolve('../../dist'),
    themeConfig: {
        logo: '/logo.png',
        lastUpdated: '最后更新于',
        nav: [
            {
                text: 'Git',
                link: 'https://github.com/tencentad/martech',
            },
        ],
        sidebar: [
            {
                title: '测试目录 - 1',
                path: '/',
                collapsable: false,
                sidebarDepth: 1,
                children: ['/common_problems'],
            },
            {
                title: '测试目录 - 2',
                path: '/',
                collapsable: false,
                sidebarDepth: 1,
                children: ['/folder-b/'],
            },
        ],
    },
    configureWebpack: (config, isServer) => {
        if (!isServer) {
            return {
                plugins: [
                    new FilemanagerWebpackPlugin({
                        events: {
                            onStart: {
                                delete: [
                                    {
                                        source: resolve('../../../cmd/web/docs'),
                                        options: {
                                            force: true,
                                        },
                                    },
                                ],
                            },
                        },
                    }),
                ],
            };
        }
    },
    async generated(pagePaths) {
        fse.copy(resolve('../../dist'), resolve('../../../cmd/web/docs'))
            .then(() => {
                console.log('Copy dist to cmd/web/docs success');
            })
            .catch((err) => {
                console.error(err);
            });
    },
};
