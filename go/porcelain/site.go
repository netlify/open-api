package porcelain

import (
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

// GetSite returns a site.
func (n *Netlify) GetSite(ctx context.Context, siteID string) (*models.Site, error) {
	authInfo := context.GetAuthInfo(ctx)
	resp, err := n.Netlify.Operations.GetSite(operations.NewGetSiteParams().WithSiteID(siteID), authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
