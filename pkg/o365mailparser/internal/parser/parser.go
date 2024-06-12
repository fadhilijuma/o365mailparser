// Package parser handles parsing of emails and extraction of attachments.
package parser

import (
	"encoding/base64"
	"fmt"
	"o365mailparser/internal/domain"
	"strings"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// ParseEmail parses a message and extracts relevant information including attachments.
func ParseEmail(message models.Messageable) (domain.Email, error) {
	email := domain.Email{
		MessageID: *message.GetId(),
	}

	// Extract attachments
	attachments, err := extractAttachments(message)
	if err != nil {
		return email, fmt.Errorf("extracting attachments: %w", err)
	}

	email.Attachments = attachments

	return email, nil
}

// extractAttachments extracts attachments from a message.
func extractAttachments(message models.Messageable) ([]domain.Attachment, error) {
	var attachments []domain.Attachment

	attachmentsCollectionPage, err := message.Attachments().Request().Get()
	if err != nil {
		return nil, fmt.Errorf("getting attachments: %w", err)
	}

	for _, attachment := range attachmentsCollectionPage.GetValue() {
		att := domain.Attachment{
			Name: *attachment.GetName(),
		}

		switch att := attachment.(type) {
		case *models.FileAttachment:
			att.Data = base64.StdEncoding.EncodeToString(att.GetContentBytes())
		case *models.ItemAttachment:
			// If it's an EML file, process it
			if strings.HasSuffix(att.GetName(), ".eml") {
				innerAttachments, err := processEmlAttachment(att)
				if err != nil {
					att.err = fmt.Errorf("processing EML attachment: %w", err)
				} else {
					attachments = append(attachments, innerAttachments...)
				}
				continue
			}
		default:
			att.err = fmt.Errorf("unsupported attachment type")
		}

		attachments = append(attachments, att)
	}

	return attachments, nil
}

// processEmlAttachment processes an EML attachment and extracts its attachments.
func processEmlAttachment(attachment *models.ItemAttachment) ([]domain.Attachment, error) {
	var attachments []domain.Attachment

	// Assuming ItemAttachment contains a message
	if message := attachment.GetItem().(models.Messageable); message != nil {
		// Retrieve the EML data from the message
		emlData := []byte(*message.GetBody().GetContent())
		innerAttachments, err := DecodeEML(emlData)
		if err != nil {
			return nil, fmt.Errorf("decoding EML: %w", err)
		}
		attachments = append(attachments, innerAttachments...)
	}

	return attachments, nil
}
