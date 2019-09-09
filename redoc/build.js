const { execFile } = require('child_process')
const { promisify } = require('util')
const { resolve } = require('path')
const cpFile = require('cp-file')

const SWAGGER_PATH = resolve(__dirname, '..', 'swagger.yml')
const OUTPUT_PATH = resolve(__dirname, 'dist', 'index.html')
const LOGO_PATH = resolve(__dirname, 'netlify-logo.png')
const OUTPUT_LOGO_PATH = resolve(__dirname, 'dist', 'netlify-logo.png')

const pExecFile = promisify(execFile)

// Build API documentation single self-contained HTML file using `redoc-cli`
const buildDocs = async function() {
  await Promise.all([redocCli(), copyAssets()])
}

const redocCli = async function() {
  await pExecFile('redoc-cli', [
    `--title=${TITLE}`,
    `--options.expandResponses=${SUCCESS_STATUS_CODES}`,
    '--options.requiredPropsFirst',
    '--options.jsonSampleExpandLevel=1',
    `--options.theme.colors.primary.main=${HEADINGS_TEXT_COLOR}`,
    `--options.theme.menu.textColor=${MENU_TEXT_COLOR}`,
    `--options.theme.menu.backgroundColor=${MENU_BACKGROUND_COLOR}`,
    `--options.theme.typography.headings.fontFamily=${FONT}`,
    `--output=${OUTPUT_PATH}`,
    'bundle',
    SWAGGER_PATH,
  ])
}

const copyAssets = async function() {
  await cpFile(LOGO_PATH, OUTPUT_LOGO_PATH)
}

const TITLE = 'Netlify API documentation'
const SUCCESS_STATUS_CODES = [200, 201, 204].join(',')
const HEADINGS_TEXT_COLOR = '#00c2b2'
const MENU_TEXT_COLOR = '#8b8b8b'
const MENU_BACKGROUND_COLOR = '#ffffff'
const FONT = 'Roboto, sans-serif'

buildDocs()
