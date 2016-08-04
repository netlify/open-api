package porcelain

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

// CustomTLSCertificate holds information
// about custom TLS certificates.
type CustomTLSCertificate struct {
	Certificate    string
	Key            string
	CACertificates string
}

// ListSites lists the sites a user has access to.
func (n *Netlify) ListSites(ctx context.Context, params *operations.ListSitesParams) ([]*models.Site, error) {
	resp, err := n.Netlify.Operations.ListSites(params, context.GetAuthWriter(ctx))
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// GetSite returns a site.
func (n *Netlify) GetSite(ctx context.Context, siteID string) (*models.Site, error) {
	authInfo := context.GetAuthWriter(ctx)
	resp, err := n.Netlify.Operations.GetSite(operations.NewGetSiteParams().WithSiteID(siteID), authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// CreateSite creates a new site.
func (n *Netlify) CreateSite(ctx context.Context, site *models.Site, configureDNS bool) (*models.Site, error) {
	authInfo := context.GetAuthWriter(ctx)

	params := operations.NewCreateSiteParams().WithSite(site).WithConfigureDNS(&configureDNS)
	resp, err := n.Netlify.Operations.CreateSite(params, authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// UpdateSite modifies an existent site.
func (n *Netlify) UpdateSite(ctx context.Context, site *models.Site) error {
	authInfo := context.GetAuthWriter(ctx)

	params := operations.NewUpdateSiteParams().WithSite(site).WithSiteID(site.ID)
	_, err := n.Netlify.Operations.UpdateSite(params, authInfo)
	return err
}

// ConfigureSiteTLS provisions a TLS certificate for a site with a custom domain.
// It uses Let's Encrypt if the certificate is empty.
func (n *Netlify) ConfigureSiteTLSCertificate(ctx context.Context, siteID string, cert *CustomTLSCertificate) (*models.SniCertificate, error) {
	authInfo := context.GetAuthWriter(ctx)

	params := operations.NewProvisionSiteTLSCertificateParams().WithSiteID(siteID)
	if cert != nil {
		params = params.WithCertificate(&cert.Certificate).WithKey(&cert.Key).WithCaCertificates(&cert.CACertificates)
	}
	resp, err := n.Netlify.Operations.ProvisionSiteTLSCertificate(params, authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// GetSiteTLSCertificate shows the TLS certificate configured for a site.
func (n *Netlify) GetSiteTLSCertificate(ctx context.Context, siteID string) (*models.SniCertificate, error) {
	authInfo := context.GetAuthWriter(ctx)

	params := operations.NewShowSiteTLSCertificateParams().WithSiteID(siteID)
	resp, err := n.Netlify.Operations.ShowSiteTLSCertificate(params, authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// WaitUntilTLSCertificateReady checks the state of a site's certificate.
// It waits until the state is "issued", for Let's Encrypt certificates
// or "custom", which means that the certificate was provided by the user.
func (n *Netlify) WaitUntilTLSCertificateReady(ctx context.Context, siteID string, cert *models.SniCertificate) (*models.SniCertificate, error) {
	if cert != nil && (cert.State == "issued" || cert.State == "custom") {
		return cert, nil
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute

	err := backoff.Retry(func() error {
		var err error
		cert, err = n.GetSiteTLSCertificate(ctx, siteID)
		if err != nil {
			return err
		}

		if cert.State != "issued" && cert.State != "custom" {
			return fmt.Errorf("certificate for site %s is not ready yet", siteID)
		}

		return nil
	}, b)

	return cert, err
}
