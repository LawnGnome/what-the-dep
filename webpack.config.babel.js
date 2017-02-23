import { resolve } from "path";
import * as webpack from "webpack";
import CopyWebpackPlugin from "copy-webpack-plugin";

export default {
  entry: "./index.js",
  output: {
    path: resolve(__dirname, "dist"),
    publicPath: "/pages/aharvey/php-crystal-ball/",
    filename: "pres.js",
  },
  devtool: "inline-source-map",
  devServer: {
    inline: true,
    outputPath: resolve(__dirname, "dist"),
  },
  module: {
    loaders: [
      {
        test: /\.css$/,
        loader: ["style-loader", "css-loader"],
      },
      {
        test: /\.less$/,
        loader: ["style-loader", "css-loader", "less-loader"],
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader",
        query: {
          presets: ["es2015"],
        },
      },
      {
        test: /\.woff2?(\?v=\d+\.\d+\.\d+)?$/,
        loader: 'url-loader',
        options: {
          limit: 50000,
        },
      },
      {
        test: /\.(png|jpg|jpeg|gif|svg)$/,
        loader: "url-loader",
        options: {
          limit: 100000,
        },
      },
      {
        test: /\.html$/,
        loader: "html-loader",
      },
    ],
  },
  plugins: [
    new webpack.optimize.OccurrenceOrderPlugin,
    new webpack.NoEmitOnErrorsPlugin,
    // Use copy-webpack-plugin to place unbundled assets in the right place.
    new CopyWebpackPlugin([
      // The stub index.html.
      {
        from: resolve(__dirname, "index.html"),
        to: resolve(__dirname, "dist", "index.html"),
      },
      // Various other dependencies.
      {
        from: resolve(__dirname, "images"),
        to: resolve(__dirname, "dist", "images"),
      },
      {
        from: resolve(__dirname, "node_modules", "reveal.js", "lib"),
        to: resolve(__dirname, "dist", "node_modules", "reveal.js", "lib"),
      },
      {
        from: resolve(__dirname, "node_modules", "reveal.js", "plugin"),
        to: resolve(__dirname, "dist", "node_modules", "reveal.js", "plugin"),
      },
    ]),
  ],
};
