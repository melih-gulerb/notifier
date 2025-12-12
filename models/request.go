package models

type SendEmailRequest struct {
	TemplateCode string                 `json:"templateCode" validate:"required"`
	To           string                 `json:"to" validate:"required,email"`
	Data         map[string]interface{} `json:"data" validate:"required"`
}
