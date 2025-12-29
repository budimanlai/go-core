package service

import (
	"context"

	common_usecase "github.com/budimanlai/go-core/common/domain/usecase"
	"gopkg.in/gomail.v2"
)

type SMTPMailServiceConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type SMTPMailServiceImpl struct {
	config              SMTPMailServiceConfig
	MessagingTemplateUC common_usecase.MessagingTemplateUsecase
}

func NewSMTPMailServiceImpl(config SMTPMailServiceConfig, template common_usecase.MessagingTemplateUsecase) SMTPMailService {
	return &SMTPMailServiceImpl{config: config, MessagingTemplateUC: template}
}

// Send sends an email with the specified parameters.
func (s *SMTPMailServiceImpl) Send(from, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// SendWithTemplate sends an email using a predefined template and data.
func (s *SMTPMailServiceImpl) SendWithTemplate(to, templateName string, templateData map[string]interface{}) error {
	// get rendered template
	tpl, err := s.MessagingTemplateUC.RenderTemplate(
		context.Background(),
		"email",
		templateName,
		templateData,
	)
	if err != nil {
		return err
	}

	// send email
	err = s.Send(s.config.From, to, tpl.Subject, tpl.ContentHtml)
	return nil
}
