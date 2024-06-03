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
	credentials, err := auth.NewClientSecretCredentials(ctx, cred)
	if err != nil {

	}
	token, err := credentials.NewClient()
	if err != nil {

	}
	fmt.Println("read emails complete", token)
	return nil
}
