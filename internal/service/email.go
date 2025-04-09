package service

import (
	"bytes"
	"cloud-sprint/config"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/go-mail/mail"
)

type EmailService struct {
	config config.EmailConfig
}

func NewEmailService(config config.EmailConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

type EmailData struct {
	To       string
	Subject  string
	Template string
	Data     map[string]interface{}
}

func (s *EmailService) SendEmail(data EmailData) error {
	templatePath := filepath.Join(s.config.TemplatesDir, data.Template)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data.Data); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	m := mail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail))
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())

	dialer := mail.NewDialer(s.config.SMTPHost, s.config.SMTPPort, s.config.SMTPUsername, s.config.SMTPPassword)
	dialer.Timeout = 10 * time.Second

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
