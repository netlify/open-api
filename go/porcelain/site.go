package porcelain

import (
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

// List the sites a user has access to.
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

// Create a new site.
func (n *Netlify) CreateSite(ctx context.Context, site *models.Site) (*models.Site, error) {
	authInfo := context.GetAuthInfo(ctx)

	resp, err := n.Netlify.Operations.CreateSite(operations.NewCreateSiteParams().WithSite(site), authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
