const { mkdtemp, rm } = require('fs/promises')
const { tmpdir } = require('os')
const { join, normalize } = require('path')

const execa = require('execa')

const { copyAssets } = require('./assets')
const { injectContent } = require('./inject')

const OUTPUT_DIR = `${__dirname}/../../dist`
const OUTPUT_PATH = `${OUTPUT_DIR}/index.html`

// Build API documentation single self-contained HTML file using `redoc-cli`
const buildDocs = async function () {
  await Promise.all([redocCli(), copyAssets(OUTPUT_DIR)])
}

// build-docs doesn't run decorators, so bundle docs@latest (redocly.yml) to a
// per-build temp file first, like `npm run bundle` does for external.yml.
const redocCli = async function () {
  const bundleDir = await mkdtemp(join(tmpdir(), 'netlify-open-api-docs-'))
  const bundlePath = join(bundleDir, 'docs.yml')
  try {
    await execa(
      'redocly',
      ['bundle', 'docs@latest', `--output=${bundlePath}`],
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
        bundlePath,
      ],
      { stdio: 'inherit' }
    )
  } finally {
    await rm(bundleDir, { recursive: true, force: true })
  }
  await injectContent(OUTPUT_PATH)
}

const TITLE = 'Netlify API documentation'
const HEADINGS_TEXT_COLOR = '#00c2b2'
const MENU_BACKGROUND_COLOR = '#ffffff'
const FONT = 'Roboto, sans-serif'
const LOGO_PADDING = '15px'

buildDocs()
