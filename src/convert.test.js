const test = require('ava')
const isPlainObj = require('is-plain-obj')

const { version } = require('../package.json')

// eslint-disable-next-line node/no-missing-require
const openApiDef = require('..')

test('OpenAPI definition normalization', async t => {
  // Ensure the OpenAPI definition general shape looks normal
  t.true(isPlainObj(openApiDef))
  t.true(isPlainObj(openApiDef.definitions))
  t.true(isPlainObj(openApiDef.paths))

  // Ensure the endpoints are present by checking for one of them
  t.true(isPlainObj(openApiDef.paths['/accounts']))
  t.true(isPlainObj(openApiDef.paths['/accounts'].get))
  t.is(openApiDef.paths['/accounts'].get.operationId, 'listAccountsForUser')

  // Ensure the host URL is present
  t.is(typeof openApiDef.host, 'string')

  // Ensure the API version is correct
  t.is(openApiDef.info.version, version)

  // Ensure JSON references are normalized
  t.true(isPlainObj(openApiDef.paths['/accounts'].get.responses.default.schema))
})
