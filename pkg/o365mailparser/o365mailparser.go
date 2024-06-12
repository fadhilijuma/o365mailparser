package o365mailparser

import (
	"context"
	"fmt"
	"o365mailparser/internal/domain"
	"o365mailparser/internal/logger"
	"o365mailparser/pkg/o365mailparser/internal/auth"
	"o365mailparser/pkg/o365mailparser/internal/email"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// FetchAndProcessEmails initializes and starts the email fetching and processing.
func FetchAndProcessEmails(log *logger.Logger, credentials domain.Credentials) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals gracefully
	shutdownChan := make(chan struct{})
	var wg sync.WaitGroup
	setupSignalHandler(cancel, shutdownChan, &wg, log)

	// Initialize authenticator
	authenticator, err := auth.NewClientSecretCredentials(ctx, credentials)
	if err != nil {
		return fmt.Errorf("creating authenticator: %w", err)
	}

	// Initialize email service
	emailService := email.NewService(ctx, authenticator)

	// Channel for receiving processed emails
	emailChan := make(chan domain.Email)

	// Start fetching and processing emails
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := emailService.FetchEmails(credentials, emailChan, shutdownChan)
		if err != nil {
			log.Error(ctx, "fetching emails", "error", err)
			cancel()
		}
	}()

	// Handle received emails
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case email := <-emailChan:
				// Process the email (e.g., save attachments, log details)
				logEmailDetails(email, log)
			case <-shutdownChan:
				return
			}
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()
	return nil
}

// setupSignalHandler sets up a signal handler for graceful shutdown.
func setupSignalHandler(cancel context.CancelFunc, shutdownChan chan struct{}, wg *sync.WaitGroup, log *logger.Logger) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals
		log.Info(context.Background(), "shutdown signal received", "signal", sig.String())
		close(shutdownChan)
		cancel()
		wg.Wait()
		os.Exit(0)
	}()
}

// logEmailDetails logs the details of the processed email.
func logEmailDetails(email domain.Email, log *logger.Logger) {
	ctx := context.Background()
	log.Info(ctx, "processed email", "MessageID", email.MessageID, "Attachments", len(email.Attachments))
	for _, attachment := range email.Attachments {
		log.Info(ctx, "attachment", "Name", attachment.Name, "Error", attachment.Data)
	}
}
