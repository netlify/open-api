const fs = require('fs')
const path = require('path')
const fromString = require('from2-string')
const pump = require('pump')
const split = require('split2')
const pt = require('parallel-transform')

const yamlPath = path.resolve(__dirname, '..', 'swagger.yml')

const swaggerYaml = fs.readFileSync(yamlPath, 'utf8')
const yamlStream = fromString(swaggerYaml)
const writeStream = fs.createWriteStream(yamlPath)

const newVersion = require('../package.json').version
let oldVersion
let didBumpVersion = false

const splitString = 'version:'

function transform (chunk, cb) {
  if (!didBumpVersion) {
    if (chunk.includes(splitString)) {
       const parts = chunk.split(splitString)
       oldVersion = parts[1]
       parts[1] = ' ' + newVersion
       chunk = parts.join(splitString)
       didBumpVersion = true
    }
  }

  chunk = chunk + '\n'
  return cb(null, chunk)
}

pump(yamlStream, split(), pt(5, transform), writeStream, (err) => {
  if (err) throw err
  if (!didBumpVersion) throw new Error('error bumping swagger.yaml version: couldnt find version string ')
  console.log('bumped swagger.yml from ' + oldVersion + ' to ' + newVersion)
})
