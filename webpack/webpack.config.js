/**
 * Created by 扬 on 2016/7/19.
 */
var webpack = require('webpack');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var CssSourcemapPlugin = require('css-sourcemaps-webpack-plugin');
module.exports = {
    entry: {
      app:  './js/app.js',
      ngapp:'./js/ngapp.js'
    },
    output: {
        path: '../assets',
        filename: '/[name].bundle.js',
        chunkFilename: "/[name].chunk.js",
        publicPath:"/"
    },
    devtool: "source-map",
    module: {
        loaders: [
            {test: /\.html$/, loader: 'raw'},
            // Transpile any JavaScript file:
            {
                test: /\.js$/,
                exclude: /node_modules/,
                loader: 'babel-loader'
            }, {
                test: /\.css$/,
                loader: ExtractTextPlugin.extract("style-loader", "css-loader")
            }, {
                test: /\.(png|jpg|gif)$/,
                loader: 'url-loader?limit=10000&name=images/[name].[ext]'
            }, {
                test   : /\.(ttf|eot|svg|woff[2]?)(\?.*)?$/,
                loader : 'file-loader?name=fonts/[name].[ext]'
            }/*,{
                test:  /\.html$/,
                loader : 'file-loader?name=templates/[name].[ext]'
            }*/
        ]
    },
    watch: true,
    watchOptions:{
        aggregateTimeout:300,
        poll:true
    },
    plugins: [
        new ExtractTextPlugin("[name].css"),
        new CssSourcemapPlugin()
    ],
    externals:{
        "jquery":"jQuery",
        "bootstrap":"bootstarp",
        "jqueryui":"jQuery"
    }
};