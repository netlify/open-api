const { tmpdir } = require('os')
const { join, normalize } = require('path')

const execa = require('execa')

const { copyAssets } = require('./assets')
const { injectContent } = require('./inject')

const OUTPUT_DIR = `${__dirname}/../../dist`
// Docs-only bundle of docs@latest (redocly.yml). build-docs doesn't run
// decorators itself, so bundle first, like `npm run bundle` does for external.yml.
const DOCS_BUNDLE_PATH = join(tmpdir(), 'netlify-open-api-docs.yml')
const OUTPUT_PATH = `${OUTPUT_DIR}/index.html`

// Build API documentation single self-contained HTML file using `redoc-cli`
const buildDocs = async function () {
  await Promise.all([redocCli(), copyAssets(OUTPUT_DIR)])
}

const redocCli = async function () {
  await execa(
    'redocly',
    ['bundle', 'docs@latest', `--output=${DOCS_BUNDLE_PATH}`],
    { stdio: 'inherit' }
  )
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
      DOCS_BUNDLE_PATH,
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
