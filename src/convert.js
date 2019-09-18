const { writeFile } = require('fs')
const { promisify } = require('util')

const SwaggerParser = require('swagger-parser')
const makeDir = require('make-dir')

const pWriteFile = promisify(writeFile)

const YAML_INPUT = `${__dirname}/../swagger.yml`
const OUTPUT_DIR = `${__dirname}/../dist`
const JSON_OUTPUT = `${OUTPUT_DIR}/swagger.json`

// Validate `swagger.yml`, dereference the JSON references then serialize to
// `swagger.json`
const convertOpenApi = async function() {
  const [openapiDef] = await Promise.all([
    SwaggerParser.validate(YAML_INPUT, { dereference: { circular: false } }),
    makeDir(OUTPUT_DIR)
  ])
  const openapiJson = JSON.stringify(openapiDef, null, 2)
  await pWriteFile(JSON_OUTPUT, openapiJson)
}

convertOpenApi()
