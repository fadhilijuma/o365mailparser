// This program provides the command line interface for the 0365mailparser service.
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/ardanlabs/conf/v3"
	"io"
	"o365mailparser/cmd/commands"
	"o365mailparser/internal/domain"
	"o365mailparser/internal/logger"
	"os"
)

var build = "develop"

type config struct {
	conf.Version
	Args        conf.Args
	Credentials struct {
		ClientID       string `conf:"default:12345666,mask"`
		ClientSecret   string `conf:"default:secret,mask"`
		TenantDomain   string `conf:"default:domain"`
		TenantID       string `conf:"default:tenant"`
		NumberOfEmails int32  `conf:"default:2"`
	}
}

func main() {
	log := logger.New(io.Discard, logger.LevelInfo, "o365mailparser", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	if err := run(log); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			fmt.Println("msg", err)
		}
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	cfg := config{
		Version: conf.Version{
			Build: build,
		},
	}

	const prefix = "o365mailparser"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		out, err := conf.String(&cfg)
		if err != nil {
			return fmt.Errorf("generating config for output: %w", err)
		}
		log.Info(context.Background(), "startup", "config", out)

		return fmt.Errorf("parsing config: %w", err)
	}

	return processCommands(cfg.Args, log, cfg)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, log *logger.Logger, cfg config) error {
	credentials := domain.Credentials{
		ClientID:       cfg.Credentials.ClientID,
		ClientSecret:   cfg.Credentials.ClientSecret,
		TenantDomain:   cfg.Credentials.TenantDomain,
		TenantID:       cfg.Credentials.TenantID,
		NumberOfEmails: cfg.Credentials.NumberOfEmails,
	}

	switch args.Num(0) {
	case "fetch":
		if err := commands.FetchAndProcessEmailsCmd(log, credentials); err != nil {
			return fmt.Errorf("fetching and processing emails: %w", err)
		}

	default:
		fmt.Println("read:       read emails from an account")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}
