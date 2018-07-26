module.exports = (ctx) => ({
  map: { inline: false },
  plugins: {
    'autoprefixer': {},
    'postcss-import': { root: ctx.file.dirname },
    'postcss-url': {
      url: 'copy',
      useHash: true,
      assetsPath: 'assets'
    },
    'postcss-browser-reporter': {},
    'postcss-reporter': {}
  }
})
