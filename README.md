# Netlify open-api

This repository contains Netlify's API definition in the [Open API format](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md), formerly known as Swagger.

It's still a work in progress and we welcome feedback and contributions.

## Go Client

We use [go-swagger](https://github.com/go-swagger/go-swagger) to validate our spec against the 2.0 spec of Open API.

We currently depend on version 0.12.0 of the swagger toolchain. You can download the binary for your platform from this release page:

https://github.com/go-swagger/go-swagger/releases/tag/0.12.0

### Spec validation

You can run this command to validate the spec:

	make validate

### Code generation

Currently, we're generating client code for Go, but we're planning on releasing libraries in any language that can generate code from the spec.

You can use this command to generate the Go client:

	make generate
	
You may first want to edit swagger.yml to add your field or endpoint definitions.

## JS Client

Swagger's JS client can dynamically create a client from a spec either from a URL or spec object.

See [swagger-js](https://github.com/swagger-api/swagger-js)


### Usage
```js
<script src='browser/swagger-client.js' type='text/javascript'></script>
<script>
var swaggerClient = new SwaggerClient('https://open-api.netlify.com/swagger.json');
</script>
```

You can also consume the swagger spec as an npm module:

```console
$ npm install @netlify/open-api
# or
$ yarn add @netlify/open-api
```

```js
import spec from '@netlify/open-api' // import the spec object into your project
```

The module also ships a copy of the original `yml` spec file at `@netlify/open-api/swagger.yml`.

## Explore API

Go to https://open-api.netlify.com to explore the spec definitions using Open-Api's UI.

## License

MIT. See [LICENSE](LICENSE) for more details.
