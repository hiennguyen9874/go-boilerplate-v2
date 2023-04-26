package emailTemplates

import (
	"context"
	"fmt"

	"github.com/matcornic/hermes/v2"
)

func (etg *emailTemplatesGenerator) GenerateVerificationCodeTemplate(ctx context.Context, name string, verificationLink string) (string, string, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				fmt.Sprintf("Welcome to %s! We're very excited to have you on board.", etg.cfg.Email.Name),
			},
			Actions: []hermes.Action{
				{
					Instructions: fmt.Sprintf("To get started with %s, please click here:", etg.cfg.Email.Name),
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Confirm your account",
						Link:  verificationLink,
					},
				},
			},
			Outros: []string{
				"Need help, or have questions? Just reply to this email, we'd love to help.",
			},
		},
	}

	// Generate an HTML email with the provided contents (for modern clients)
	emailBody, err := etg.h.GenerateHTML(email)
	if err != nil {
		return "", "", nil
	}

	// Generate the plaintext version of the e-mail (for clients that do not support xHTML)
	emailText, err := etg.h.GeneratePlainText(email)
	if err != nil {
		return "", "", nil
	}

	return emailBody, emailText, nil
}
