const { readFile, writeFile } = require('fs')
const { promisify } = require('util')

const { version } = require('../package.json')

const pReadFile = promisify(readFile)
const pWriteFile = promisify(writeFile)

const SWAGGER_PATH = `${__dirname}/../swagger.yml`

// Modify the `info.version` field inside `swagger.yml` so it matches the
// `package.json` `version`
const bumpSwagger = async function() {
  const swaggerFile = await pReadFile(SWAGGER_PATH, 'utf8')
  const newSwaggerFile = swaggerFile.replace(VERSION_REGEXP, `$1${version}`)

  if (newSwaggerFile !== swaggerFile) {
    await pWriteFile(SWAGGER_PATH, newSwaggerFile)
  }
}

const VERSION_REGEXP = /^(\s+version: )([\d.]+)$/m

bumpSwagger()
