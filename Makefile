.PONY: all build generate help validate

all: validate generate build ## Validate the swagger spec, generate the code and build it.

build: ## Build the API Go client.
	go build go/...

generate: validate ## Generate the API Go client and the JSON document for the UI.
	swagger generate client -A netlify -f swagger.yml -t go -c plumbing
	swagger generate spec -i swagger.yml -o ui/swagger.json

help: ## Show this help.
		@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


validate: ## Check that the swagger spec is valid.
	swagger validate swagger.yml
