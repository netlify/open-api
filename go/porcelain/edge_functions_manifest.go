package porcelain

// https://github.com/netlify/build/blob/main/packages/edge-bundler/node/manifest.ts
type edgeFunctionsManifest struct {
	Bundles []edgeFunctionsManifestBundle `json:"bundles"`
}

type edgeFunctionsManifestBundle struct {
	Asset  string `json:"asset"`
	Format string `json:"format"`
}
