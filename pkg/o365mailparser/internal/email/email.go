// Package email manages email retrieval and basic email operations.
package email

import (
	"context"
	"fmt"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"o365mailparser/internal/domain"
	"o365mailparser/pkg/o365mailparser/internal/auth"
	"o365mailparser/pkg/o365mailparser/internal/parser"
)

type Service struct {
	ctx            context.Context
	clientProvider auth.ClientProvider
}

// NewService creates a new instance of Service.
func NewService(ctx context.Context, authenticator auth.ClientProvider) *Service {
	return &Service{
		ctx:            ctx,
		clientProvider: authenticator,
	}
}

// FetchEmails fetches emails from the Office 365 mailbox.
func (es *Service) FetchEmails(creds domain.Credentials, emailChan chan<- domain.Email, shutdownChan <-chan struct{}) error {
	client, err := es.clientProvider.NewClient()
	if err != nil {
		return fmt.Errorf("creating graph client: %w", err)
	}

	// Get messages from the mailbox
	userID := creds.TenantDomain // Assuming the user ID is the tenant domain
	messagesRequest := client.UsersById(userID).Messages().Request()
	messagesRequest.Top(int(creds.NumberOfEmails))

	mailResponse, err := messagesRequest.Get(es.ctx)
	if err != nil {
		return fmt.Errorf("getting messages: %w", err)
	}

	// Process each message
	for _, message := range mailResponse.GetValue() {
		select {
		case <-shutdownChan:
			return nil
		default:
			email, err := es.processMessage(message)
			if err != nil {
				fmt.Printf("error processing message: %v\n", err)
				continue
			}
			emailChan <- email
		}
	}

	return nil
}

// processMessage processes an individual message and extracts attachments and metadata.
func (es *Service) processMessage(message models.Messageable) (domain.Email, error) {
	return parser.ParseEmail(message)
}
