package porcelain

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
	"github.com/netlify/open-api/v2/go/porcelain/context"

	"github.com/sirupsen/logrus"
)

type SiteAsset struct {
	SiteID  string
	Name    string
	Size    int64
	Private bool
	Body    io.ReadSeeker
}

func (n *Netlify) UploadNewSiteAsset(ctx context.Context, asset *SiteAsset) (*models.Asset, error) {
	buffer := make([]byte, 512)
	b, err := asset.Body.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	contentType := http.DetectContentType(buffer[:b])
	params := operations.NewCreateSiteAssetParams().WithSiteID(asset.SiteID).
		WithName(asset.Name).WithSize(asset.Size).WithContentType(contentType)

	if asset.Private {
		visibility := "private"
		params = params.WithVisibility(&visibility)
	}

	signature, err := n.AddSiteAsset(ctx, params)
	if err != nil {
		return nil, err
	}

	asset.Body.Seek(0, 0)
	bufferWriter := &bytes.Buffer{}
	writer := multipart.NewWriter(bufferWriter)

	for key, value := range signature.Form.Fields {
		writer.WriteField(key, value)
	}

	part, err := writer.CreateFormFile("file", asset.Name)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, asset.Body); err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", signature.Form.URL, bufferWriter)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		defer resp.Body.Close()
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected error uploading assets to the :cloud: %d - %s", resp.StatusCode, respBody)
	}

	updateParams := operations.NewUpdateSiteAssetParams().WithSiteID(asset.SiteID).WithAssetID(signature.Asset.ID).WithState("uploaded")

	ar, err := n.UpdateSiteAsset(ctx, updateParams)
	if err != nil {
		return nil, err
	}

	return ar, nil
}

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

func (n *Netlify) ListSiteAssets(ctx context.Context, params *operations.ListSiteAssetsParams) ([]*models.Asset, error) {
	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id": params.SiteID,
	}).Debug("Listing site assets")

	resp, err := n.Netlify.Operations.ListSiteAssets(params, context.GetAuthInfo(ctx))
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

func (n *Netlify) ShowSiteAssetInfo(ctx context.Context, params *operations.GetSiteAssetInfoParams, showSignature bool) (*models.Asset, error) {
	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id":  params.SiteID,
		"asset_id": params.AssetID,
	}).Debug("Show site asset information")

	authInfo := context.GetAuthInfo(ctx)

	resp, err := n.Netlify.Operations.GetSiteAssetInfo(params, authInfo)
	if err != nil {
		return nil, err
	}

	asset := resp.Payload
	if asset.Visibility == "private" && showSignature {
		sigParams := operations.NewGetSiteAssetPublicSignatureParams().WithSiteID(params.SiteID).WithAssetID(params.AssetID)

		sig, err := n.GetSiteAssetPublicSignature(ctx, sigParams)
		if err != nil {
			return nil, err
		}

		asset.URL = sig.URL
	}

	return asset, nil
}

func (n *Netlify) GetSiteAssetPublicSignature(ctx context.Context, params *operations.GetSiteAssetPublicSignatureParams) (*models.AssetPublicSignature, error) {
	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id":  params.SiteID,
		"asset_id": params.AssetID,
	}).Debug("Get site asset public signature")

	authInfo := context.GetAuthInfo(ctx)

	resp, err := n.Netlify.Operations.GetSiteAssetPublicSignature(params, authInfo)
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}
