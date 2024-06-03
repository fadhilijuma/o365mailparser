package domain

// Email holds parsing of one email.
type Email struct {
	MessageID   string       // Message ID of the email
	Attachments []Attachment // List of attachments in the email
}

// Attachment holds parsing of one attachment.
type Attachment struct {
	Name string // File Name of the attachment
	Data string // Base64 encoded data of the attachment
	err  error  // Error in parsing the individual attachment
}

// Credentials are used by OAuthClient.
type Credentials struct {
	ClientID       string
	ClientSecret   string
	TenantDomain   string
	TenantID       string
	NumberOfEmails int
}
