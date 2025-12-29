package service

type WhatsAppService interface {
	// SendMessage sends a WhatsApp message to the specified recipient.
	SendMessage(to string, message string) error

	// SendMessageWithTemplate sends a WhatsApp message using a predefined template.
	SendWithTemplate(to string, templateName string, data map[string]interface{}) error
}
