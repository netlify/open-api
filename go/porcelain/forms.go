package porcelain

import (
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

// ListForms lists the forms a user has access to.
func (n *Netlify) ListForms(ctx context.Context, params *operations.ListFormsParams) ([]*models.Form, error) {
	resp, err := n.Netlify.Operations.ListForms(params, context.GetAuthInfo(ctx))
	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

// ListFormsBySiteId lists the forms of a particular site
func (n *Netlify) ListFormsBySiteId(ctx context.Context, siteID string) ([]*models.Form, error) {
	authInfo := context.GetAuthInfo(ctx)
	resp, err := n.Netlify.Operations.ListForms(operations.NewListFormsParams().WithSiteID(&siteID), authInfo)
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
