package sendEmail

import (
	"context"
	"crypto/tls"

	"github.com/hiennguyen9874/go-boilerplate-v2/config"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(ctx context.Context, from string, to string, subject string, bodyHtml string, bodyPlain string) error
}

type emailSender struct {
	cfg *config.Config
}

func NewEmailSender(cfg *config.Config) EmailSender {
	return &emailSender{
		cfg: cfg,
	}
}

func (es *emailSender) SendEmail(ctx context.Context, from string, to string, subject string, bodyHtml string, bodyPlain string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyHtml)
	m.AddAlternative("text/plain", bodyPlain)

	d := gomail.NewDialer(es.cfg.SmtpEmail.Host, es.cfg.SmtpEmail.Port, es.cfg.SmtpEmail.User, es.cfg.SmtpEmail.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
