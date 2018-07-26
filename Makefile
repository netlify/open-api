.PHONY: all build deps generate help test validate
CHECK_FILES?=$$(go list ./... | grep -v /vendor/)




help: ## Show this help.
		@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: validate generate build ## Validate the swagger spec, generate the code and build it.

build: ## Build the API Go client.
	cd ./go && go build ./...

deps: ## Download dependencies.
	go get -u github.com/golang/dep/cmd/dep && cd ./go && dep ensure

generate: validate ## Generate the API Go client and the JSON document for the UI.
	swagger generate client -A netlify -f swagger.yml -t go -c plumbing

test: ## Test the go code.
	cd ./go && go test -v $(CHECK_FILES)

validate: ## Check that the swagger spec is valid.
	swagger validate swagger.yml
