package porcelain

import (
	"github.com/go-openapi/runtime"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
)

// GetSite returns a site.
func (n *Netlify) GetSite(siteID string, authInfo runtime.ClientAuthInfoWriter) (*models.Site, error) {
	resp, err := n.Netlify.Operations.GetSite(operations.NewGetSiteParams().WithSiteID(siteID), authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
