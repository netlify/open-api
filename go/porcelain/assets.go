package porcelain

import (
	"github.com/Sirupsen/logrus"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

func (n *Netlify) AddSiteAsset(ctx context.Context, params *operations.CreateSiteAssetParams) (*models.AssetSignature, error) {
	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id": params.SiteID,
	}).Debug("Creating site asset signature")

	resp, err := n.Netlify.Operations.CreateSiteAsset(params, context.GetAuthInfo(ctx))
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

func (n *Netlify) UpdateSiteAsset(ctx context.Context, params *operations.UpdateSiteAssetParams) (*models.Asset, error) {
	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id": params.SiteID,
	}).Debug("Updating site asset state")

	resp, err := n.Netlify.Operations.UpdateSiteAsset(params, context.GetAuthInfo(ctx))
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}
