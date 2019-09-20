const { readFile, writeFile } = require('fs')
const { promisify } = require('util')

const HEAD_PATH = `${__dirname}/head.html`

const pReadFile = promisify(readFile)
const pWriteFile = promisify(writeFile)

// Inject HTML content after Redoc has built the documentation
const injectContent = async function(outputPath) {
  const [siteContent, head] = await Promise.all([pReadFile(outputPath, 'utf8'), pReadFile(HEAD_PATH, 'utf8')])

  const updatedContent = siteContent.replace(END_HEAD_REGEXP, `${head}$&`)

  await pWriteFile(outputPath, updatedContent)
}

const END_HEAD_REGEXP = /<\/head>/

module.exports = { injectContent }
