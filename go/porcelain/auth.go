package porcelain

import (
	"fmt"
	"time"

	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
	"github.com/netlify/open-api/v2/go/porcelain/context"

	"github.com/sirupsen/logrus"
)

const (
	ticketingTimeout = time.Minute * 5
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

func (n *Netlify) WaitUntilTicketAuthorized(ctx context.Context, ticket *models.Ticket) (*models.Ticket, error) {
	authInfo := context.GetAuthInfo(ctx)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	params := operations.NewShowTicketParams().WithTicketID(ticket.ID)
	start := time.Now()
	for t := range ticker.C {
		resp, err := n.Netlify.Operations.ShowTicket(params, authInfo)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		context.GetLogger(ctx).WithFields(logrus.Fields{
			"ticket_id":  ticket.ID,
			"authorized": resp.Payload.Authorized,
		}).Debug("Waiting until deploy ready")

		if resp.Payload.Authorized {
			return resp.Payload, nil
		}

		if t.Sub(start) > ticketingTimeout {
			return nil, fmt.Errorf("Error: the authorization process timed out")
		}
	}

	return ticket, nil
}
