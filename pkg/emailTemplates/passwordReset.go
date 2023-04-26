package emailTemplates

import (
	"context"
	"fmt"

	"github.com/matcornic/hermes/v2"
)

func (etg *emailTemplatesGenerator) GeneratePasswordResetTemplate(ctx context.Context, name string, resetLink string) (string, string, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				fmt.Sprintf("You have received this email because a password reset request for %s account was received.", etg.cfg.Email.Name),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password (valid for 15min):",
					Button: hermes.Button{
						Color: "#DC4D2F", // Optional action button color
						Text:  "Reset your password",
						Link:  resetLink,
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
			Signature: "Thanks",
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
