package porcelain

import (
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

// Create a login ticket to authenticate a user
func (n *Netlify) CreateTicket(ctx context.Context, clientID string) (*models.Ticket, error) {
	params := operations.NewCreateTicketParams().WithClientID(clientID)
	resp, err := n.Netlify.Operations.CreateTicket(params, context.GetAuthInfo(ctx))

	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

func (n *Netlify) ShowTicket(ctx context.Context, ticketID string) (*models.Ticket, error) {
	params := operations.NewShowTicketParams().WithTicketID(ticketID)
	resp, err := n.Netlify.Operations.ShowTicket(params, context.GetAuthInfo(ctx))

	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}

func (n *Netlify) ExchangeTicket(ctx context.Context, ticketID string) (*models.AccessToken, error) {
	params := operations.NewExchangeTicketParams().WithTicketID(ticketID)
	resp, err := n.Netlify.Operations.ExchangeTicket(params, context.GetAuthInfo(ctx))

	if err != nil {
		return nil, err
	}

	return resp.Payload, nil
}
