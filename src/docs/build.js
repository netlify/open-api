const { normalize } = require('path')

const execa = require('execa')

const { copyAssets } = require('./assets')
const { injectContent } = require('./inject')

const SWAGGER_PATH = `${__dirname}/../../external.yml`
const OUTPUT_DIR = `${__dirname}/../../dist`
const OUTPUT_PATH = `${OUTPUT_DIR}/index.html`

// Build API documentation single self-contained HTML file using `redoc-cli`
const buildDocs = async function () {
  await Promise.all([redocCli(), copyAssets(OUTPUT_DIR)])
}

const redocCli = async function () {
  await execa(
    'redocly',
    [
      `--title=${TITLE}`,
      '--theme.openapi.requiredPropsFirst',
      '--theme.openapi.sortOperationsAlphabetically',
      `--theme.openapi.colors.primary.main=${HEADINGS_TEXT_COLOR}`,
      `--theme.openapi.sidebar.backgroundColor=${MENU_BACKGROUND_COLOR}`,
      `--theme.openapi.typography.headings.fontFamily=${FONT}`,
      `--theme.openapi.logo.gutter=${LOGO_PADDING}`,
      `--output=${normalize(OUTPUT_PATH)}`,
      'build-docs',
      SWAGGER_PATH,
    ],
    { stdio: 'inherit' }
  )
  await injectContent(OUTPUT_PATH)
}

const TITLE = 'Netlify API documentation'
const HEADINGS_TEXT_COLOR = '#00c2b2'
const MENU_BACKGROUND_COLOR = '#ffffff'
const FONT = 'Roboto, sans-serif'
const LOGO_PADDING = '15px'

buildDocs()
