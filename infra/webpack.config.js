const path = require('path');

module.exports = {
  entry: './handler.js',
  target: 'node',
  mode: 'production',
  optimization: {
    minimize: false,
  },
  performance: {
    hints: false,
  },
  devtool: 'source-map',
  externals: ['aws-sdk'],
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env'],
          },
        },
      },
    ],
  },
  resolve: {
    extensions: ['.js', '.json'],
  },
  output: {
    libraryTarget: 'commonjs2',
    path: path.resolve(__dirname, '.webpack'),
    filename: '[name].js',
  },
};