// Package process contains the main logic for processing emails, such as downloading emails and extracting information.
package process

import (
	"context"
	"log"
	"o365mailparser/internal/domain"
	"o365mailparser/pkg/o365mailparser/internal/email"
)

type Processor struct {
	emailService *email.Service
	emailChan    chan domain.Email
	shutdownChan chan struct{}
	resultChan   chan domain.Email
}

// NewProcessor creates a new instance of Processor.
func NewProcessor(ctx context.Context, emailService *email.Service, resultChan chan domain.Email) *Processor {
	return &Processor{
		emailService: emailService,
		emailChan:    make(chan domain.Email),
		shutdownChan: make(chan struct{}),
		resultChan:   resultChan,
	}
}

// Start begins processing emails.
func (p *Processor) Start(creds domain.Credentials) {
	go func() {
		if err := p.emailService.FetchEmails(creds, p.emailChan, p.shutdownChan); err != nil {
			log.Fatalf("Failed to fetch emails: %v", err)
		}
		close(p.emailChan)
	}()

	go p.processEmails()
}

// Stop signals the processor to stop.
func (p *Processor) Stop() {
	close(p.shutdownChan)
}

// processEmails processes emails received from the email channel.
func (p *Processor) processEmails() {
	for email := range p.emailChan {
		// Perform any additional processing if needed.
		p.resultChan <- email
	}
	close(p.resultChan)
}
