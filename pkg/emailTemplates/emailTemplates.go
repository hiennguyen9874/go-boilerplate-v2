package emailTemplates

import (
	"context"

	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/matcornic/hermes/v2"
)

type EmailTemplatesGenerator interface {
	GenerateVerificationCodeTemplate(ctx context.Context, name string, verificationLink string) (string, string, error)
	GeneratePasswordResetTemplate(ctx context.Context, name string, resetLink string) (string, string, error)
}

type emailTemplatesGenerator struct {
	h   *hermes.Hermes
	cfg *config.Config
}

func NewEmailTemplatesGenerator(cfg *config.Config) EmailTemplatesGenerator {
	return &emailTemplatesGenerator{
		h: &hermes.Hermes{
			// Optional Theme
			Theme: new(hermes.Default),
			Product: hermes.Product{
				// Appears in header & footer of e-mails
				Name: cfg.Email.Name,
				Link: cfg.Email.Link,
				// Optional product logo
				Logo:      cfg.Email.LogoLink,
				Copyright: cfg.Email.Copyright,
			},
		},
		cfg: cfg,
	}
}
