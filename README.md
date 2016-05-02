# Introduction

This repository contains Netlify's API definition in the [Open API format](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md), formerly known as Swagger.

It's still a work in progress and we welcome feedback and contributions.

## Validation

We use [go-swagger](https://github.com/go-swagger/go-swagger) to validate our spec against the 2.0 spec of Open API.

Follow the installation instruction in that repository. You can run this command to validate the spec:

	swagger validate -f swagger.yml

## Code generation

Currently, we're generating client code for Go, but we're planning on releasing libraries in any language that can generate code from the spec.

You can use this command to generate the Go client:

	swagger generate client -t go -f swagger.yml

## License

MIT. See [LICENSE](LICENSE) for more details.
