all: validate generate

generate:
	swagger generate client -A netlify -f swagger.yml -t go -c plumbing

validate:
	swagger validate swagger.yml
