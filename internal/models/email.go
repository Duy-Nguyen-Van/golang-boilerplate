package models

import (
	"time"
)

// EmailMessage represents an email message domain entity
type EmailMessage struct {
	ID         string     `json:"id" db:"id"`
	To         []string   `json:"to" db:"to"`
	CC         []string   `json:"cc,omitempty" db:"cc"`
	BCC        []string   `json:"bcc,omitempty" db:"bcc"`
	Subject    string     `json:"subject" db:"subject"`
	Body       string     `json:"body" db:"body"`
	HTMLBody   string     `json:"html_body,omitempty" db:"html_body"`
	Status     string     `json:"status" db:"status"` // pending, sent, failed
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	SentAt     *time.Time `json:"sent_at,omitempty" db:"sent_at"`
	ErrorMsg   string     `json:"error_msg,omitempty" db:"error_msg"`
	RetryCount int        `json:"retry_count" db:"retry_count"`
}

// IsValid checks if the email message is valid
func (e *EmailMessage) IsValid() bool {
	return len(e.To) > 0 && e.Subject != "" && (e.Body != "" || e.HTMLBody != "")
}

// EmailTemplate represents an email template
type EmailTemplate struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Subject   string    `json:"subject" db:"subject"`
	Body      string    `json:"body" db:"body"`
	HTMLBody  string    `json:"html_body" db:"html_body"`
	Variables []string  `json:"variables" db:"variables"` // JSON array of variable names
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IsValid checks if the email template is valid
func (t *EmailTemplate) IsValid() bool {
	return t.Name != "" && t.Subject != "" && (t.Body != "" || t.HTMLBody != "")
}
