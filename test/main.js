import test from 'ava'

import openApiDef from '..'

test('OpenAPI definition snapshot', async t => {
  t.snapshot(openApiDef)
})
