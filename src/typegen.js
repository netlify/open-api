const fs = require('fs')

async function typegen() {
  const openapiTS = (await import('openapi-typescript')).default
  const openapiJsonString = fs.readFileSync('dist/openapi.json', 'utf-8')
  const openapi = JSON.parse(openapiJsonString)

  addDocsLinksToOperationDescriptions(openapi)

  const result = await openapiTS(openapi, {
    xNullableAsNullable: true,
  })

  fs.writeFileSync('dist/index.d.ts', result)
}

typegen()

/**
 * Add a link to the Netlify OpenAPI reference docs for each operation in the description.
 * This is used in the generated TypeScript types to provide a link to the OpenAPI docs on each operation.
 */
function addDocsLinksToOperationDescriptions(openapi) {
  for (const path of Object.keys(openapi.paths)) {
    for (const method of Object.keys(openapi.paths[path])) {
      const tag =
        openapi.paths[path][method].tags && openapi.paths[path][method].tags[0]

      const operationId = openapi.paths[path][method].operationId

      if (operationId && tag) {
        const description = openapi.paths[path][method].description

        const docsLink = `API Reference: {@link https://open-api.netlify.com/#tag/${tag}/operation/${operationId} | \`${operationId}\`} `

        openapi.paths[path][method].description = `${docsLink}\n\n${
          description ? `${description}\n\n` : ''
        }`
      }
    }
  }
}
