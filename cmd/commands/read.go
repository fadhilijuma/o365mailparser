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

// Read reads the number of emails in a specified email account.
func Read(log *logger.Logger, credentials domain.Credentials) error {
	o365mailparser.Read(log, credentials)
	fmt.Println("read emails complete", credentials)
	return nil
}
