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

## Node Client

In order to generate the Node client, you must install Java 8 and the swagger-codegen tool.

1. You will need java8 (or 7) to run the code-gen

    ```bash
    # install java 8 via homebrew
    brew cask install homebrew/cask-versions/java8
    ```


2. Then install swagger codegen

    ```bash
    # install code-gen via homebrew
    brew install swagger-codegen
    ```

3. Generate the library

	```console
	$ make node
	```

The `swagger-codegen` generates an entire module in the `node` folder.  You can publish new versions of the client from that folder with `npm publish`.

## Explore API

Go to https://open-api.netlify.com to explore the spec definitions using Open-Api's UI.

## License

MIT. See [LICENSE](LICENSE) for more details.
