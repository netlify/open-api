# CONTRIBUTING

Contributions are welcome!

The go-client is an [netlify/open-api][open-api] derived http client generated using [go-swagger][go-swagger]. It uses go modules. To work on it please ensure the following:

- You are running at least Go 1.12
- You have cloned this repo OUTSIDE of the go path. (So go modules work).
- You have a $GOPATH set up and $GOPATH/bin is added to your \$PATH

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

1. Make sure your PR title and commits follow the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) spec.
2. Don't bump the version number for `swagger.yml` changes. The release process handles that.
3. Ensure `make validate` passes.
4. The go tests run against the last generated go client. These must pass before making a release.
5. If all you want is a new endpoint, you can PR just the `swagger.yml` changes for review and regenerate the go client when its ready to go in.

## Making a new release

Merge the release PR (auto generated via `release-please`)

## License

By contributing to Netlify Node Client, you agree that your contributions will be licensed
under its [MIT license](LICENSE).

[godoc-img]: https://godoc.org/github.com/netlify/go-client/?status.svg
[godoc]: https://godoc.org/github.com/netlify/go-client
[goreport-img]: https://goreportcard.com/badge/github.com/netlify/go-client
[goreport]: https://goreportcard.com/report/github.com/netlify/go-client
[git-img]: https://img.shields.io/github/release/netlify/go-client.svg
[git]: https://github.com/netlify/go-client/releases/latest
[modules]: https://github.com/golang/go/wiki/Modules
[open-api]: https://github.com/netlify/open-api
[go-swagger]: https://github.com/go-swagger/go-swagger
[go-modules]: https://github.com/golang/go/wiki/Modules
[swagger]: https://github.com/netlify/open-api/blob/master/swagger.yml
