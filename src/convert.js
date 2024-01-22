const { writeFile, mkdir } = require('fs/promises')

const SwaggerParser = require('swagger-parser')

const YAML_INPUT = `${__dirname}/../swagger.yml`
const OUTPUT_DIR = `${__dirname}/../dist`
const JSON_OUTPUT = `${OUTPUT_DIR}/swagger.json`

// Validate `swagger.yml`, dereference the JSON references then serialize to
// `swagger.json`
const convertOpenApi = async function () {
  const [openapiDef] = await Promise.all([
    SwaggerParser.validate(YAML_INPUT, { dereference: { circular: false } }),
    mkdir(OUTPUT_DIR, { recursive: true }),
  ])
  const openapiJson = JSON.stringify(openapiDef, null, 2)
  await writeFile(JSON_OUTPUT, openapiJson)
}

convertOpenApi()
