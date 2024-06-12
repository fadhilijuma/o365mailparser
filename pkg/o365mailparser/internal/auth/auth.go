package auth

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	graph "github.com/microsoftgraph/msgraph-sdk-go"
	"o365mailparser/internal/domain"
)

// https://learn.microsoft.com/en-us/office/office-365-management-api/get-started-with-office-365-management-apis

type CredentialProvider interface {
	NewCredentials() (*azidentity.ClientSecretCredential, error)
}

type ClientProvider interface {
	NewClient(cred *azidentity.ClientSecretCredential) (*graph.GraphServiceClient, error)
}

type Auth struct {
	credProvider   CredentialProvider
	clientProvider ClientProvider
}

func NewAuth(credProvider CredentialProvider, clientProvider ClientProvider) *Auth {
	return &Auth{
		credProvider:   credProvider,
		clientProvider: clientProvider,
	}
}

type ClientSecretCredentialProvider struct {
	ctx   context.Context
	creds domain.Credentials
}

func NewClientSecretCredentialProvider(ctx context.Context, c domain.Credentials) *ClientSecretCredentialProvider {
	return &ClientSecretCredentialProvider{ctx: ctx, creds: c}
}

// NewCredentials returns an authenticated confidential client using the provided tenant credentials.
func (c *ClientSecretCredentialProvider) NewCredentials() (*azidentity.ClientSecretCredential, error) {
	cred, err := azidentity.NewClientSecretCredential(c.creds.TenantID, c.creds.ClientID, c.creds.ClientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("creating credentials: %w", err)
	}
	return cred, nil
}

// NewClient returns to us a brand-new token or the existing token from the cache that is not yet expired.
func (a *Auth) NewClient() (*graph.GraphServiceClient, error) {
	cred, err := a.credProvider.NewCredentials()
	if err != nil {
		return nil, fmt.Errorf("failed to create credentials: %w", err)
	}

	client, err := a.clientProvider.NewClient(cred)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}
