package porcelain

import (
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
	"github.com/netlify/open-api/v2/go/porcelain/context"
)

// ListFormsBySiteId lists the forms of a particular site
func (n *Netlify) ListFormsBySiteId(ctx context.Context, siteID string) ([]*models.Form, error) {
	authInfo := context.GetAuthInfo(ctx)
	resp, err := n.Netlify.Operations.ListSiteForms(operations.NewListSiteFormsParams().WithSiteID(siteID), authInfo)
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// ListFormSubmissions lists the forms submissions of a particular form
func (n *Netlify) ListFormSubmissions(ctx context.Context, formID string) ([]*models.Submission, error) {
	authInfo := context.GetAuthInfo(ctx)
	resp, err := n.Netlify.Operations.ListFormSubmissions(operations.NewListFormSubmissionsParams().WithFormID(formID), authInfo)
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}
