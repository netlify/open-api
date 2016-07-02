package porcelain

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/netlify/open-api/go/plumbing"
)

// Default netlify HTTP client.
var Default = NewHTTPClient(nil)

// NewHTTPClient creates a new netlify HTTP client.
func NewHTTPClient(formats strfmt.Registry) *Netlify {
	n := plumbing.NewHTTPClient(formats)
	return &Netlify{n}
}

// New creates a new netlify client
func New(transport runtime.ClientTransport, formats strfmt.Registry) *Netlify {
	n := plumbing.New(transport, formats)
	return &Netlify{n}
}

// Netlify is a client for netlify
type Netlify struct {
	*plumbing.Netlify
}
