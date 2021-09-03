const webpack = require('webpack');
const path = require('path')
const json5 = require('json5');
const {VueLoaderPlugin} = require('vue-loader');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = (env, args) => {
  return {
    entry: './src/main.js',
    output: {
      filename: '[name].[contenthash:8].bundle.js',
      path: path.resolve(__dirname, './dist'),
      clean: true,
    },
    module: {
      rules: [
        {
          test: /\.vue$/,
          loader: 'vue-loader',
        },
        {
          test: /\.s[ac]ss$/,
          use: [
            'vue-style-loader',
            'css-loader',
            'sass-loader',
          ],
        },
        {
          test: /\.css$/,
          use: [
            'vue-style-loader',
            'css-loader',
          ],
        },
        {
          test: /\.json5$/,
          type: 'json',
          parser: {
            parse: json5.parse,
          },
        }
      ],
    },
    plugins: [
      new VueLoaderPlugin(),
      new HtmlWebpackPlugin({
        filename: 'index.html',
        template: 'index.html',
      }),
      new webpack.DefinePlugin({
        'PRODUCTION': args.mode == 'production',
      }),
    ],
    devServer: {
      static: './dist',
    },
    devtool: (args.mode == 'production' ? 'source-map' : 'inline-source-map'),
  };
};
