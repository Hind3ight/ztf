/* eslint-disable @typescript-eslint/no-var-requires */
const bodyParser = require('body-parser')
const { NODE_ENV, VUE_APP_PORT, VUE_APP_MOCK } = process.env;

const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin')

module.exports = {
    publicPath: '/',
    outputDir: 'dist',
    productionSourceMap: false,
    devServer: {
        host: 'localhost',
        port: VUE_APP_PORT || 8000,
        disableHostCheck: true,
        // 配置反向代理
        /*
        proxy: {
            '/api': {
              target: '<url>',
              ws: true,
              changeOrigin: true
            },
            '/foo': {
              target: '<other_url>'
            }
        },
        */
        before: function(app, server) {
        }
    },
    css: {
        loaderOptions: {
            less: {
                javascriptEnabled: true,
            }
        }
    },

    transpileDependencies: [],
    configureWebpack: {
        // 不需要打包的插件
        externals: {
            // 'vue': 'Vue',
            // 'vue-router': 'VueRouter',
        },
    },

    chainWebpack(config) {
        // 内置的 svg Rule 添加 exclude
        config.module
        .rule('svg')
        .exclude.add(/iconsvg/)
        .end();

        // 添加 svg-sprite-loader Rule
        config.module
        .rule('svg-sprite-loader')
        .test(/.svg$/)
        .include.add(/iconsvg/)
        .end()
        .use('svg-sprite-loader')
        .loader('svg-sprite-loader');

        // 添加 svgo Rule
        config.module
        .rule('svgo')
        .test(/.svg$/)
        .include.add(/iconsvg/)
        .end()
        .use('svgo-loader')
        .loader('svgo-loader')
        .options({
            // externalConfig 配置特殊不是相对路径，起始路径是根目录
            externalConfig: './src/assets/iconsvg/svgo.yml',
        });

        config.resolve.alias.set('vue-i18n', 'vue-i18n/dist/vue-i18n.cjs.js')

        config.plugin('monaco-editor').use(MonacoWebpackPlugin, [
            {
                // Languages are loaded on demand at runtime
                languages: ['javascript', 'php']
            }
        ])
    }
}
