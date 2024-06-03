package o365mailparser

import (
	"context"
	"fmt"
	"o365mailparser/internal/domain"
	"o365mailparser/internal/logger"
	"o365mailparser/pkg/o365mailparser/internal/auth"
)

// Read reads the number of emails in a specified email account.
func Read(ctx context.Context, log *logger.Logger, cred domain.Credentials) error {
	authenticate, err := auth.Authenticate(ctx, cred)
	if err != nil {

	}
	token, err := authenticate.AcquireToken()
	if err != nil {

	}
	fmt.Println("read emails complete", token)
	return nil
}
