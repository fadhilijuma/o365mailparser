package auth

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	"o365mailparser/internal/domain"
)

// https://learn.microsoft.com/en-us/office/office-365-management-api/get-started-with-office-365-management-apis

type Authenticator interface {
	NewClient() (*graph.GraphServiceClient, error)
}

// Auth wraps a contextualised confidential client.
type Auth struct {
	ctx  context.Context
	cred *azidentity.ClientSecretCredential
}

// NewClientSecretCredentials returns an authenticated confidential client using the provided tenant credentials.
func NewClientSecretCredentials(ctx context.Context, c domain.Credentials) (*Auth, error) {
	cred, err := azidentity.NewClientSecretCredential(c.TenantID, c.ClientID, c.ClientSecret, nil)

	if err != nil {
		return nil, fmt.Errorf("creating credentials: %w", err)
	}
	return &Auth{ctx: ctx, cred: cred}, nil
}

// NewClient returns to us a brand-new token or the existing token from the cache that is not yet expired.
func (a *Auth) NewClient() (*graph.GraphServiceClient, error) {
	client, err := graph.NewGraphServiceClientWithCredentials(a.cred, []string{"Files.Read"})
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}
	return client, nil
}
