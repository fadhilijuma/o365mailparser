package auth

import (
	"context"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"o365mailparser/internal/domain"
)

// https://learn.microsoft.com/en-us/office/office-365-management-api/get-started-with-office-365-management-apis
// https://login.windows.net/common/oauth2/authorize?response_type=code&resource=https%3A%2F%2Fmanage.office.com&client_id={your_client_id}&redirect_uri={your_redirect_url }
var (
	defaultApplicationID = "o365mailparser"

	microsoftTokenURL = "https://login.microsoftonline.com/%s"
)

type Auth struct {
	ctx    context.Context
	client confidential.Client
}

// Authenticate returns an authenticated confidential client using the provided tenant credentials.
func Authenticate(ctx context.Context, c *domain.Credentials) (*Auth, error) {
	cred, err := confidential.NewCredFromSecret(c.ClientSecret)
	if err != nil {
		return nil, err
	}
	confidentialClient, err := confidential.New(fmt.Sprintf(microsoftTokenURL, c.TenantDomain), defaultApplicationID, cred)
	if err != nil {
		return nil, fmt.Errorf("authenticate tenant credentiats: %w", err)
	}
	return &Auth{ctx: ctx, client: confidentialClient}, nil
}

// AcquireToken returns to us a brand-new token or the existing token from the cache that is not yet expired.
func (a *Auth) AcquireToken() (string, error) {
	scopes := []string{"email"}
	result, err := a.client.AcquireTokenSilent(a.ctx, scopes)
	if err != nil {
		// cache miss, authenticate with another AcquireToken... method
		result, err = a.client.AcquireTokenByCredential(a.ctx, scopes)
		if err != nil {
			return "", fmt.Errorf("acquire token silently: %w", err)
		}
	}
	accessToken := result.AccessToken
	return accessToken, nil
}
