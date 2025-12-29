package service

type SMTPMailService interface {
	// Send sends an email with the specified parameters.
	Send(from, to, subject, body string) error

	// SendWithTemplate adds an email sending job to the queue using a template.
	SendWithTemplate(to, templateName string, templateData map[string]interface{}) error
}
