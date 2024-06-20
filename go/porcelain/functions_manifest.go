package porcelain

import "github.com/netlify/open-api/v2/go/models"

// https://github.com/netlify/zip-it-and-ship-it/blob/main/src/manifest.ts
type functionsManifest struct {
	Functions []functionsManifestEntry `json:"functions"`
	Version   int                      `json:"version"`
}

type functionsManifestEntry struct {
	MainFile       string                  `json:"mainFile"`
	Name           string                  `json:"name"`
	Path           string                  `json:"path"`
	Runtime        string                  `json:"runtime"`
	RuntimeVersion string                  `json:"runtimeVersion"`
	Schedule       string                  `json:"schedule"`
	DisplayName    string                  `json:"displayName"`
	Generator      string                  `json:"generator"`
	Timeout        int64                   `json:"timeout"`
	BuildData      map[string]interface{}  `json:"buildData"`
	InvocationMode string                  `json:"invocationMode"`
	Routes         []functionRoute         `json:"routes"`
	ExcludedRoutes []excludedFunctionRoute `json:"excludedRoutes"`
	Priority       int                     `json:"priority"`
	TrafficRules   *functionTrafficRules   `json:"trafficRules"`
}

type functionRoute struct {
	Pattern      string   `json:"pattern"`
	Literal      string   `json:"literal"`
	Expression   string   `json:"expression"`
	Methods      []string `json:"methods"`
	PreferStatic bool     `json:"prefer_static"`
}

type excludedFunctionRoute struct {
	Pattern    string `json:"pattern"`
	Literal    string `json:"literal"`
	Expression string `json:"expression"`
}

type functionTrafficRules struct {
	Action struct {
		Type   string `json:"type"`
		Config struct {
			RateLimitConfig struct {
				Algorithm   string `json:"algorithm"`
				WindowSize  int    `json:"windowSize"`
				WindowLimit int    `json:"windowLimit"`
			}
			Aggregate *models.TrafficRulesAggregateConfig
			To        string `json:"to"`
		}
	}
}
