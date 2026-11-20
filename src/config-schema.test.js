const { readFile } = require('fs/promises')

const test = require('ava')

const { getOutputPath } = require('./config-schema.js')

test('netlify.toml JSON schema is published at its $id', async (t) => {
  const schemaJson = await readFile(
    require.resolve('@netlify/config/schema.json'),
    'utf8'
  )
  const schema = JSON.parse(schemaJson)

  t.is(await readFile(getOutputPath(schema), 'utf8'), schemaJson)
})

test('netlify.toml JSON schema must have an $id on this site', (t) => {
  t.throws(() =>
    getOutputPath({ $id: 'https://example.com/netlify.toml.json' })
  )
  t.throws(() =>
    getOutputPath({
      $id: 'https://open-api.netlify.com/schemas/netlify.toml.json',
    })
  )
})
