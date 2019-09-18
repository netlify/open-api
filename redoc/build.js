const { execFile } = require('child_process')
const { promisify } = require('util')
const { resolve, join } = require('path')
const cpFile = require('cp-file')
const fs = require('fs')
const readFile = promisify(fs.readFile)
const writeFile = promisify(fs.writeFile)

const SWAGGER_PATH = resolve(__dirname, '..', 'swagger.yml')
const LOGO_PATH = resolve(__dirname, 'netlify-logo.png')
const OUTPUT_DIR = resolve(__dirname, 'dist')
const OUTPUT_PATH = resolve(OUTPUT_DIR, 'index.html')
const OUTPUT_JS_PATH = resolve(OUTPUT_DIR, 'js')
const OUTPUT_LOGO_PATH = resolve(OUTPUT_DIR, 'netlify-logo.png')

const pExecFile = promisify(execFile)

// Build API documentation single self-contained HTML file using `redoc-cli`
const buildDocs = async function() {
  await Promise.all([redocCli(), copyAssets()])
  // Inject analytics
  await addAnalytics()
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
    `--output=${OUTPUT_PATH}`,
    'bundle',
    SWAGGER_PATH
  ])
}

const TITLE = 'Netlify API documentation'
const HEADINGS_TEXT_COLOR = '#00c2b2'
const MENU_TEXT_COLOR = '#8b8b8b'
const MENU_BACKGROUND_COLOR = '#ffffff'
const FONT = 'Roboto, sans-serif'
const LOGO_PADDING = '15px'
const GOOGLE_ANALYTICS = 'UA-42258181-19'

const copyAssets = async function() {
  await cpFile(LOGO_PATH, OUTPUT_LOGO_PATH)
}

const ANALYTICS_SCRIPT_PATH = resolve('node_modules/analytics/dist/analytics.min.js')
const GA_SCRIPT_PATH = resolve('node_modules/analytics-plugin-ga/dist/analytics-plugin-ga.min.js')
const addAnalytics = async function() {
  // Add scripts to dist
  await cpFile(ANALYTICS_SCRIPT_PATH, join(OUTPUT_JS_PATH, 'analytics.min.js'))
  await cpFile(GA_SCRIPT_PATH, join(OUTPUT_JS_PATH, 'analytics-plugin-ga.min.js'))
  // Inject JS into html
  const siteContent = await readFile(OUTPUT_PATH, 'utf-8')
  const analyticsScript = `
  <!-- Include analytics -->
  <script src="/js/analytics.min.js"></script>
  <script src="/js/analytics-plugin-ga.min.js"></script>
  <!-- initialize analytics -->
  <script type="text/javascript">
    /* Initialize analytics */
    var Analytics = _analytics.init({
      plugins: [
        analyticsGA({
          trackingId: '${GOOGLE_ANALYTICS}'
        })
      ]
    })
    Analytics.page()
    </script>
  </head>
  `
  const updatedContent = siteContent.replace(/<\/head>/g, analyticsScript)
  await writeFile(OUTPUT_PATH, updatedContent)
}

buildDocs()
