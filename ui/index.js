const SwaggerUI = require('swagger-ui')

SwaggerUI({
  dom_id: '#container',
  url: '/swagger.json',
  displayOperationId: true,
  deepLinking: true
})
