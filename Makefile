all: validate generate

generate: validate
	swagger generate client -A netlify -f swagger.yml -t go -c plumbing
	swagger generate spec -i swagger.yml -o ui/swagger.json

validate:
	swagger validate swagger.yml
