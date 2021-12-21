package porcelain

// https://github.com/netlify/zip-it-and-ship-it/blob/main/src/manifest.ts
type functionsManifest struct {
	Functions []functionsManifestEntry `json:"functions"`
	Version   int                      `json:"version"`
}

type functionsManifestEntry struct {
	MainFile string `json:"mainFile"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Runtime  string `json:"runtime"`
	Schedule string `json:"schedule"`
}
