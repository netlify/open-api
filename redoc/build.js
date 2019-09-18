const { readFile, writeFile } = require('fs')
const { execFile } = require('child_process')
const { promisify } = require('util')
const { basename, normalize } = require('path')

const cpFile = require('cp-file')

const SWAGGER_PATH = `${__dirname}/../swagger.yml`
const OUTPUT_DIR = `${__dirname}/dist`
const OUTPUT_PATH = `${OUTPUT_DIR}/index.html`
const HEAD_PATH = `${__dirname}/head.html`

const ASSETS = [
  `${__dirname}/netlify-logo.png`,
  `${__dirname}/../node_modules/analytics/dist/analytics.min.js`,
  `${__dirname}/../node_modules/analytics-plugin-ga/dist/analytics-plugin-ga.min.js`
]

const pExecFile = promisify(execFile)
const pReadFile = promisify(readFile)
const pWriteFile = promisify(writeFile)

// Build API documentation single self-contained HTML file using `redoc-cli`
const buildDocs = async function() {
  await Promise.all([redocCli(), copyAssets()])
}

const redocCli = async function() {
  await pExecFile('redoc-cli', [
    `--title=${TITLE}`,
    '--options.requiredPropsFirst',
    `--options.theme.colors.primary.main=${HEADINGS_TEXT_COLOR}`,
    `--options.theme.menu.textColor=${MENU_TEXT_COLOR}`,
    `--options.theme.menu.backgroundColor=${MENU_BACKGROUND_COLOR}`,
    `--options.theme.typography.headings.fontFamily=${FONT}`,
    `--options.theme.logo.gutter=${LOGO_PADDING}`,
    `--output=${normalize(OUTPUT_PATH)}`,
    'bundle',
    SWAGGER_PATH
  ])
  await injectContent()
}

const TITLE = 'Netlify API documentation'
const HEADINGS_TEXT_COLOR = '#00c2b2'
const MENU_TEXT_COLOR = '#8b8b8b'
const MENU_BACKGROUND_COLOR = '#ffffff'
const FONT = 'Roboto, sans-serif'
const LOGO_PADDING = '15px'

// Inject HTML content after Redoc has built the documentation
const injectContent = async function() {
  const [siteContent, head] = await Promise.all([pReadFile(OUTPUT_PATH, 'utf8'), pReadFile(HEAD_PATH, 'utf8')])

  const updatedContent = siteContent.replace(END_HEAD_REGEXP, head)

  await pWriteFile(OUTPUT_PATH, updatedContent)
}

const END_HEAD_REGEXP = /<\/head>/

// Copy files for static site to use (logo, JavaScript libraries)
const copyAssets = async function() {
  await Promise.all(ASSETS.map(copyAsset))
}

const copyAsset = async function(path) {
  await cpFile(path, `${OUTPUT_DIR}/${basename(path)}`)
}

buildDocs()
