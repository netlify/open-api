# CONTRIBUTING

We use [go-swagger](https://github.com/go-swagger/go-swagger) to validate our spec against the 2.0 spec of Open API.

We currently depend on version 0.18.0 of the go swagger toolchain. You can download the binary for your platform from this release page:

https://github.com/go-swagger/go-swagger/releases/tag/0.18.0

## Spec validation

You can run this command to validate the spec:

	make validate

Always validate after making changes to the `swagger.yml` file.

### Go Client generation

The Go client must be regenerated after every change to the `swagger.yml`.

You can use this command to generate the Go client:

	make generate
	
You may first want to edit swagger.yml to add your field or endpoint definitions.

## Making PRs

1. Don't bump the version number for `swagger.yml` changes.  Do that during the release process.
2. Ensure `make validate` passes.
3. The go tests run against the last generated go client.  These must pass before making a release.
4. If all you want is a new endpoint, you can PR just the `swagger.yml` changes for review and regenerate the go client when its ready to go in.

## Making a new release

1. bump the version of swagger.yml file (after making changes to it)
2. regenarate go client (if you haven't)
3. bump a JS package version with `npm version [major|minor|patch]` (updates package.json, create a git tag)
4. make sure everything is committed and `git push && git push --tags` to push to the origin
5. write a release note for the tag in [Releases](https://github.com/netlify/open-api/releases) page
6. publish to npm (`npm install && npm publish`)

## License

By contributing to Netlify Node Client, you agree that your contributions will be licensed
under its [MIT license](LICENSE).
