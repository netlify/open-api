const { readFile, writeFile, mkdir } = require('fs/promises')
const { basename } = require('path')

const SCHEMA_INPUT = require.resolve('@netlify/config/schema.json')
const OUTPUT_DIR = `${__dirname}/../dist`
const PUBLIC_ORIGIN = 'https://open-api.netlify.com'

// The `netlify.toml` JSON schema is maintained alongside the configuration
// validation in `@netlify/config`. We publish it on this site, at the URL that
// the schema declares as its `$id`.
const getOutputPath = function ({ $id }) {
  const url = new URL($id)

  if (
    url.origin !== PUBLIC_ORIGIN ||
    url.pathname !== `/${basename(url.pathname)}`
  ) {
    throw new Error(`The netlify.toml schema has an unexpected $id: ${$id}`)
  }

  return `${OUTPUT_DIR}${url.pathname}`
}

const copyConfigSchema = async function () {
  const [schemaJson] = await Promise.all([
    readFile(SCHEMA_INPUT, 'utf8'),
    mkdir(OUTPUT_DIR, { recursive: true }),
  ])
  await writeFile(getOutputPath(JSON.parse(schemaJson)), schemaJson)
}

module.exports = { getOutputPath }

if (require.main === module) {
  copyConfigSchema()
}
