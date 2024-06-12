package commands

import (
	"errors"
	"fmt"
	"o365mailparser/internal/domain"
	"o365mailparser/internal/logger"
	"o365mailparser/pkg/o365mailparser"
)

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

// FetchAndProcessEmailsCmd reads and processes the specified number of emails in a specified email account.
func FetchAndProcessEmailsCmd(log *logger.Logger, credentials domain.Credentials) error {
	if err := o365mailparser.FetchAndProcessEmails(log, credentials); err != nil {
		return fmt.Errorf("fetching and processing emails: %w", err)
	}
	fmt.Println("fetch and process emails complete", credentials)
	return nil
}
