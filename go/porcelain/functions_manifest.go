package porcelain

// https://github.com/netlify/zip-it-and-ship-it/blob/main/src/manifest.ts
type functionsManifest struct {
	Functions []functionsManifestEntry `json:"functions"`
	Version   int                      `json:"version"`
}

type functionsManifestEntry struct {
	MainFile       string `json:"mainFile"`
	Name           string `json:"name"`
	Path           string `json:"path"`
	Runtime        string `json:"runtime"`
	RuntimeVersion string `json:"runtimeVersion"`
	Schedule       string `json:"schedule"`
	DisplayName    string `json:"displayName"`
	Generator      string `json:"generator"`
	InvocationMode string `json:"invocation_mode"`
}
