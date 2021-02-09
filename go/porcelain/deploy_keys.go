package porcelain

import (
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
	"github.com/netlify/open-api/v2/go/porcelain/context"
)

func (n *Netlify) CreateDeployKey(ctx context.Context) (*models.DeployKey, error) {
	authInfo := context.GetAuthInfo(ctx)
	params := operations.NewCreateDeployKeyParams()
	resp, err := n.Netlify.Operations.CreateDeployKey(params, authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
