const fs = require('fs')
const path = require('path')
const yaml = require('js-yaml')
const swaggerYaml = fs.readFileSync(path.resolve(__dirname, '..', 'swagger.yml'), 'utf8')
const dfn = yaml.safeLoad(swaggerYaml)

const OUTPUT = path.join(__dirname, 'dist')


fs.writeFileSync(path.join(OUTPUT, 'swagger.json'), JSON.stringify(dfn, null, '\t'))
fs.writeFileSync(path.join(OUTPUT, 'swagger.yml'), swaggerYaml)
