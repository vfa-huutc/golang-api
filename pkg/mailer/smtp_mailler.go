package mailer

import "gopkg.in/gomail.v2"

type EmailSender interface {
	Send(to []string, subject, body string) error
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type GomailSender struct {
	config SMTPConfig
}

func NewGomailSender(config SMTPConfig) *GomailSender {
	return &GomailSender{config: config}
}

// Send sends an email using the SMTP configuration
// Parameters:
//   - to: slice of recipient email addresses
//   - subject: email subject line
//   - body: plain text email body content
//
// Returns:
//   - error if sending fails, nil on success
func (s *GomailSender) Send(to []string, subject, html string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)

	dialer := gomail.NewDialer(s.config.Host, s.config.Port, s.config.Username, s.config.Password)
	return dialer.DialAndSend()
}
