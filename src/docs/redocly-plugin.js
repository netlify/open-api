// Docs-only Redocly decorator (wired up in redocly.yml for `docs@latest`).
// Redoc uses an operation's `summary` as its heading and sidebar label, which
// hides the operationId. `netlify api --help` sends people to this page to
// find method names, so prefix each summary with its operationId in the docs.
// The published swagger.json and openapi.json keep the plain summary.
const operationIdInSummary = () => ({
  Operation: {
    leave(operation) {
      if (operation.operationId && operation.summary) {
        operation.summary = `${operation.operationId}: ${operation.summary}`
      }
    },
  },
})

module.exports = {
  id: 'netlify-docs',
  decorators: {
    oas2: { 'operation-id-in-summary': operationIdInSummary },
    oas3: { 'operation-id-in-summary': operationIdInSummary },
  },
}
