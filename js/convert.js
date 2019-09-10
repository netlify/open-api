const { writeFile } = require('fs')
const { promisify } = require('util')

const SwaggerParser = require('swagger-parser')

const pWriteFile = promisify(writeFile)

const YAML_INPUT = `${__dirname}/../swagger.yml`
const JSON_OUTPUT = `${__dirname}/dist/swagger.json`

const convertOpenApi = async function() {
  const openapiDef = await SwaggerParser.validate(YAML_INPUT, {
    dereference: { circular: false },
  })
  const openapiJson = JSON.stringify(openapiDef, null, 2)
  await pWriteFile(JSON_OUTPUT, openapiJson)
}

convertOpenApi()
