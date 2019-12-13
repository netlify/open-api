package tools

// todo: Find a better way to ignore tags or make a breaking release
//go:generate gobin -m -run github.com/go-swagger/go-swagger/cmd/swagger flatten swagger.yml -o swagger_flat.json
//go:generate sh -c "cat swagger_flat.json | jq '[., (.paths | map_values(.[] |= del(.tags?)) | {paths: .})] | add' > swagger_go.json"
//go:generate gobin -m -run github.com/go-swagger/go-swagger/cmd/swagger generate client -A netlify -f swagger_go.json -t go -c plumbing --default-scheme=https --with-flatten=full
