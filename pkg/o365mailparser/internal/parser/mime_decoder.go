package parser

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"o365mailparser/internal/domain"
	"strings"
)

// DecodeEML decodes an EML file and extracts its attachments.
func DecodeEML(emlData []byte) ([]domain.Attachment, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(emlData))
	if err != nil {
		return nil, fmt.Errorf("reading EML message: %w", err)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("parsing media type: %w", err)
	}

	var attachments []domain.Attachment

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])

		for {
			part, err := mr.NextPart()
			if err != nil {
				break
			}

			attachment, err := processPart(part)
			if err != nil {
				return nil, fmt.Errorf("processing part: %w", err)
			}
			attachments = append(attachments, attachment)
		}
	}

	return attachments, nil
}

// processPart processes a single part of a MIME message and returns it as an attachment.
func processPart(part *multipart.Part) (domain.Attachment, error) {
	var attachment domain.Attachment
	attachment.Name = part.FileName()

	partData, err := ioutil.ReadAll(part)
	if err != nil {
		return attachment, fmt.Errorf("reading part data: %w", err)
	}

	attachment.Data = base64.StdEncoding.EncodeToString(partData)

	return attachment, nil
}
