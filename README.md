# Introduction

This repository contains Netlify's API definition in the [Open API format](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/2.0.md), formerly known as Swagger.

It's still a work in progress and we welcome feedback and contributions.

## Installation

We use [go-swagger](https://github.com/go-swagger/go-swagger) to validate our spec against the 2.0 spec of Open API.

To download the toolchain follow the installation instructions in that repository. You'll need to have Go installed.

## Spec validation

You can run this command to validate the spec:

	make validate

## Code generation

Currently, we're generating client code for Go, but we're planning on releasing libraries in any language that can generate code from the spec.

You can use this command to generate the Go client:

	make generate

## Explore API

Go to https://open-api.netlify.com to explore the spec definitions using Open-Api's UI.

## License

MIT. See [LICENSE](LICENSE) for more details.
