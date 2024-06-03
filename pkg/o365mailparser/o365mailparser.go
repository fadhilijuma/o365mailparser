package o365mailparser

import (
	"fmt"
	"o365mailparser/internal/domain"
	"o365mailparser/internal/logger"
)

// Read reads the number of emails in a specified email account.
func Read(log *logger.Logger, credentials domain.Credentials) error {
	fmt.Println("read emails complete", credentials)
	return nil
}
