package porcelain

import (
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

func (n *Netlify) ListSites(ctx context.Context, params *operations.ListSitesParams) ([]*models.Site, error) {
	resp, err := n.Netlify.Operations.ListSites(params, context.GetAuthInfo(ctx))
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// GetSite returns a site.
func (n *Netlify) GetSite(ctx context.Context, siteID string) (*models.Site, error) {
	authInfo := context.GetAuthInfo(ctx)
	resp, err := n.Netlify.Operations.GetSite(operations.NewGetSiteParams().WithSiteID(siteID), authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
