const test = require('ava')

// eslint-disable-next-line node/no-missing-require
const openApiDef = require('..')

test('OpenAPI definition snapshot', async t => {
  t.snapshot(openApiDef)
})
