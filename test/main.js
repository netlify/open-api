const test = require('ava')

const openApiDef = require('..')

test('OpenAPI definition snapshot', async t => {
  t.snapshot(openApiDef)
})
