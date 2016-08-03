package porcelain

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/netlify/open-api/go/plumbing"
)

var defaultStreamingHost = "streaming.netlify.com"
var defaultStreamingPath = "/api/v1/streaming"

// Default netlify HTTP client.
var Default = NewHTTPClient(nil)

// NewHTTPClient creates a new netlify HTTP client.
func NewHTTPClient(formats strfmt.Registry) *Netlify {
	n := plumbing.NewHTTPClient(formats)
	ret := &Netlify{n, nil}
	ret.SetStreamingEndpoint("wss://" + defaultStreamingHost + defaultStreamingPath)
	return ret
}

// New creates a new netlify client
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Netlify {
	n := plumbing.New(transport, formats)
	ret := &Netlify{n, nil}
	ret.SetStreamingEndpoint("wss://" + defaultStreamingHost + defaultStreamingPath)
	return ret
}

// Netlify is a client for netlify
type Netlify struct {
	*plumbing.Netlify
	streamingURL *url.URL
}

// SetStreamingEndpoint will parse the endpoint and sanity check the URL
func (n *Netlify) SetStreamingEndpoint(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}

	if u.Scheme == "" {
		u.Scheme = "ws"
	}

	if !(u.Scheme == "ws" || u.Scheme == "wss") {
		return fmt.Errorf("Unsupported schemed '%s' - only 'ws' or 'wss' are allowed", u.Scheme)
	}

	if u.Host == "" {
		u.Host = defaultStreamingHost
	}

	if u.Path == "" {
		u.Path = defaultStreamingPath
	}

	n.streamingURL = u
	return nil
}
