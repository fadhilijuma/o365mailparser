package auth

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"o365mailparser/internal/domain"
	"time"
)

var (
	defaultBaseURL   = "https://manage.office.com"
	defaultVersion   = "v1.0"
	defaultUserAgent = "o365mailparser"
	defaultTimeout   = 5 * time.Second

	microsoftTokenURL = "https://login.windows.net/%s/oauth2/token?api-version=1.0"
)

// OAuthClient returns an authenticated httpClient using the provided credentials.
func OAuthClient(ctx context.Context, c *domain.Credentials) *http.Client {
	conf := &clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     fmt.Sprintf(microsoftTokenURL, c.TenantDomain),
		EndpointParams: url.Values{
			"resource": []string{defaultBaseURL},
		},
	}
	return conf.Client(ctx)
}
