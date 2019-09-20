const { basename } = require('path')

const cpFile = require('cp-file')

const ASSETS = [
  `${__dirname}/netlify-logo.png`,
  `${__dirname}/script.js`,
  `${__dirname}/../../node_modules/analytics/dist/analytics.min.js`,
  `${__dirname}/../../node_modules/analytics-plugin-ga/dist/analytics-plugin-ga.min.js`
]

// Copy files for static site to use (logo, JavaScript libraries)
const copyAssets = async function(outputDir) {
  await Promise.all(ASSETS.map(path => copyAsset(path, outputDir)))
}

const copyAsset = async function(path, outputDir) {
  await cpFile(path, `${outputDir}/${basename(path)}`)
}

module.exports = { copyAssets }
