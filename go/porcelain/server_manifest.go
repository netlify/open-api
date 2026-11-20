package porcelain

// serverManifest is what the server build step writes beside the archive it
// produced, describing only the deploy's Netlify Server.
//
// https://github.com/netlify/build/blob/main/packages/zip-it-and-ship-it/src/manifest.ts
type serverManifest struct {
	Server  *serverManifestEntry `json:"server"`
	Version int                  `json:"version"`
}

type serverManifestEntry struct {
	Path           string                  `json:"path"`
	Region         string                  `json:"region"`
	InstanceSize   string                  `json:"instanceSize"`
	Routes         []functionRoute         `json:"routes"`
	ExcludedRoutes []excludedFunctionRoute `json:"excludedRoutes"`
	RuntimeVersion string                  `json:"runtimeVersion"`
	Scale          *serverManifestScale    `json:"scale"`
}

type serverManifestScale struct {
	Min *int64 `json:"min"`
	Max *int64 `json:"max"`
}
