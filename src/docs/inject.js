const { readFile, writeFile } = require('fs')
const { promisify } = require('util')

const HEAD_PATH = `${__dirname}/head.html`
const FOOTER_PATH = `${__dirname}/footer.html`

const pReadFile = promisify(readFile)
const pWriteFile = promisify(writeFile)

// Inject HTML content after Redoc has built the documentation
const injectContent = async function (outputPath) {
  const [siteContent, head, footer] = await Promise.all([
    pReadFile(outputPath, 'utf8'),
    pReadFile(HEAD_PATH, 'utf8'),
    pReadFile(FOOTER_PATH, 'utf8'),
  ])

  const updatedContent = siteContent
    .replace(END_HEAD_REGEXP, `${head}$&`)
    .replace(END_BODY_REGEXP, `${footer}$&`)

  await pWriteFile(outputPath, updatedContent)
}

const END_HEAD_REGEXP = /<\/head>/
const END_BODY_REGEXP = /<\/body>/

module.exports = { injectContent }
