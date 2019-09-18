# CONTRIBUTING

Contributions are welcome!

The go-client is an [netlify/open-api][open-api] derived http client generated using [go-swagger][go-swagger]. It uses go modules. To work on it please ensure the following:

- You are running at least Go 1.12
- You have cloned this repo OUTSIDE of the go path. (So go modules work).
- You have a $GOPATH set up and $GOPATH/bin is added to your \$PATH

See [GMBE:Tools as dependencies](https://github.com/go-modules-by-example/index/tree/master/010_tools) and [GMBE:Using `gobin` to install/run tools](https://github.com/go-modules-by-example/index/tree/master/017_using_gobin) for a deeper explanation of how a tools.go file works.

## Spec validation

All spec changes must pass go-swagger spec validation.

You can run this command to validate the spec:

    make validate

Always validate after making changes to the `swagger.yml` file.

### Go Client generation

The Go client must be regenerated after every change to the `swagger.yml`.

You can use this command to generate the Go client:

    make generate

You may first want to edit swagger.yml to add your field or endpoint definitions.

## Making PRs

1. Don't bump the version number for `swagger.yml` changes. Do that during the release process.
2. Ensure `make validate` passes.
3. The go tests run against the last generated go client. These must pass before making a release.
4. If all you want is a new endpoint, you can PR just the `swagger.yml` changes for review and regenerate the go client when its ready to go in.

## Making a new release

1. Make sure you are on the HEAD of the master branch.
2. regenarate go client (if you haven't) (Make all and commit the results)
3. bump a JS package version with `npm version [major|minor|patch]` (updates package.json, swagger.yaml and create a git tag)
4. Run `npm publish` which will as `git push && git push --tags` to push to the origin, create a github release and publish the spec to npm.

## License

By contributing to Netlify Node Client, you agree that your contributions will be licensed
under its [MIT license](LICENSE).

[godoc-img]: https://godoc.org/github.com/netlify/go-client/?status.svg
[godoc]: https://godoc.org/github.com/netlify/go-client
[goreport-img]: https://goreportcard.com/badge/github.com/netlify/go-client
[goreport]: https://goreportcard.com/report/github.com/netlify/go-client
[git-img]: https://img.shields.io/github/release/netlify/go-client.svg
[git]: https://github.com/netlify/go-client/releases/latest
[gobin]: https://github.com/myitcv/gobin
[modules]: https://github.com/golang/go/wiki/Modules
[open-api]: https://github.com/netlify/open-api
[go-swagger]: https://github.com/go-swagger/go-swagger
[go-modules]: https://github.com/golang/go/wiki/Modules
[swagger]: https://github.com/netlify/open-api/blob/master/swagger.yml
