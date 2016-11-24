package porcelain

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/netlify/open-api/go/plumbing"
)

const DefaultSyncFileLimit = 7000
const DefaultConcurrentUploadLimit = 10

// Default netlify HTTP client.
var Default = NewHTTPClient(nil)

// NewHTTPClient creates a new netlify HTTP client.
func NewHTTPClient(formats strfmt.Registry) *Netlify {
	n := plumbing.NewHTTPClient(formats)
	return &Netlify{
		Netlify:       n,
		syncFileLimit: DefaultSyncFileLimit,
		uploadLimit:   DefaultConcurrentUploadLimit,
	}
}

// New creates a new netlify client
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
	n.syncFileLimit = limit
}

func (n *Netlify) SetConcurrentUploadLimit(limit int) {
	if limit > 0 {
		n.uploadLimit = limit
	}
}
