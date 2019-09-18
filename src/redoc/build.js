const { execFile } = require('child_process')
const { promisify } = require('util')
const { normalize } = require('path')

const { copyAssets } = require('./assets')
const { injectContent } = require('./inject')

const SWAGGER_PATH = `${__dirname}/../swagger.yml`
const OUTPUT_DIR = `${__dirname}/../../dist`
const OUTPUT_PATH = `${OUTPUT_DIR}/index.html`

const pExecFile = promisify(execFile)

// Build API documentation single self-contained HTML file using `redoc-cli`
const buildDocs = async function() {
  await Promise.all([redocCli(), copyAssets(OUTPUT_DIR)])
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
  await injectContent(OUTPUT_PATH)
}

const TITLE = 'Netlify API documentation'
const HEADINGS_TEXT_COLOR = '#00c2b2'
const MENU_TEXT_COLOR = '#8b8b8b'
const MENU_BACKGROUND_COLOR = '#ffffff'
const FONT = 'Roboto, sans-serif'
const LOGO_PADDING = '15px'

buildDocs()
