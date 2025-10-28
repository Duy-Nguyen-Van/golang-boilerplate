package dtos

import "time"

// EmailSendRequest represents an email send request DTO
type EmailSendRequest struct {
	To       []string `json:"to" validate:"required,min=1"`
	CC       []string `json:"cc,omitempty"`
	BCC      []string `json:"bcc,omitempty"`
	Subject  string   `json:"subject" validate:"required,min=1,max=200"`
	Body     string   `json:"body" validate:"required,min=1"`
	HTMLBody string   `json:"html_body,omitempty"`
}

// EmailTemplateRequest represents an email template request DTO
type EmailTemplateRequest struct {
	Name      string   `json:"name" validate:"required,min=1,max=100"`
	Subject   string   `json:"subject" validate:"required,min=1,max=200"`
	Body      string   `json:"body" validate:"required,min=1"`
	HTMLBody  string   `json:"html_body,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

// EmailTemplateResponse represents an email template response DTO
type EmailTemplateResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	HTMLBody  string    `json:"html_body"`
	Variables []string  `json:"variables"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EmailSendTemplateRequest represents a template-based email send request DTO
type EmailSendTemplateRequest struct {
	TemplateID string                 `json:"template_id" validate:"required"`
	To         []string               `json:"to" validate:"required,min=1"`
	CC         []string               `json:"cc,omitempty"`
	BCC        []string               `json:"bcc,omitempty"`
	Variables  map[string]interface{} `json:"variables,omitempty"`
}

// EmailResponse represents an email response DTO
type EmailResponse struct {
	ID        string     `json:"id"`
	To        []string   `json:"to"`
	Subject   string     `json:"subject"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	SentAt    *time.Time `json:"sent_at,omitempty"`
	ErrorMsg  string     `json:"error_msg,omitempty"`
}

// EmailListResponse represents a list of emails response
type EmailListResponse struct {
	Emails     []EmailResponse `json:"emails"`
	Pagination Pageable        `json:"pagination"`
}
