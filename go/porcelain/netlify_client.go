package porcelain

import (
	"github.com/netlify/open-api/v2/go/plumbing"
	"github.com/netlify/open-api/v2/go/porcelain/http"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

const DefaultSyncFileLimit = 500
const DefaultConcurrentUploadLimit = 10
const DefaultRetryAttempts = 3

// Default netlify HTTP client.
var Default = NewHTTPClient(nil)

// NewHTTPClient creates a new netlify HTTP client.
func NewHTTPClient(formats strfmt.Registry) *Netlify {
	cfg := plumbing.DefaultTransportConfig()
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)

	return New(transport, formats)
}

// NewRetryableHTTPClient creates a new netlify HTTP client with a number of attempts for rate limits.
func NewRetryableHTTPClient(formats strfmt.Registry, attempts int) *Netlify {
	cfg := plumbing.DefaultTransportConfig()
	transport := httptransport.New(cfg.Host, cfg.BasePath, cfg.Schemes)

	return NewRetryable(transport, formats, attempts)
}

// NewRetryable creates a new netlify client with a number of attempts for rate limits.
func NewRetryable(transport runtime.ClientTransport, formats strfmt.Registry, attempts int) *Netlify {
	tr := http.NewRetryableTransport(transport, attempts)
	return New(tr, formats)
}

// New creates a new netlify client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Netlify {
	n := plumbing.New(transport, formats)
	return &Netlify{
		Netlify:       n,
		syncFileLimit: DefaultSyncFileLimit,
		uploadLimit:   DefaultConcurrentUploadLimit,
	}
}

// Netlify is a client for netlify
type Netlify struct {
	*plumbing.Netlify
	syncFileLimit int
	uploadLimit   int
}

func (n *Netlify) SetSyncFileLimit(limit int) {
	if limit > 0 {
		n.syncFileLimit = limit
	}
}

func (n *Netlify) SetConcurrentUploadLimit(limit int) {
	if limit > 0 {
		n.uploadLimit = limit
	}
}
