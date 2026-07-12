// Package authmail contains email logic for auth-related emails.
package authmail

import (
	"context"
	"fmt"

	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail"
	"github.com/andrew-hayworth22/critiquefi/services/api/internal/mail/htmlrender"
)

// AuthMail sends mail related to auth
type AuthMail struct {
	sender  mail.Sender
	baseURL string
}

// New creates a new AuthMail instance.
func New(sender mail.Sender, baseURL string) *AuthMail {
	return &AuthMail{sender: sender, baseURL: baseURL}
}

// SendWelcome sends a welcome email to the user.
func (a *AuthMail) SendWelcome(ctx context.Context, recipient mail.Recipient) error {
	html, err := htmlrender.RenderWelcomeTemplate(htmlrender.WelcomeData{
		Name: recipient.Name,
	})
	if err != nil {
		return err
	}

	email := mail.Email{
		Subject: "welcome to critiquefi!",
		Body:    html,
		To:      []mail.Recipient{recipient},
	}

	return a.sender.Send(ctx, email)
}

// SendPasswordReset sends a password reset email to the user.
func (a *AuthMail) SendPasswordReset(ctx context.Context, recipient mail.Recipient, token string) error {
	html, err := htmlrender.RenderPasswordResetTemplate(htmlrender.PasswordResetData{
		Name:     recipient.Name,
		ResetURL: fmt.Sprintf("%s/reset-password?token=%s", a.baseURL, token),
	})
	if err != nil {
		return err
	}

	email := mail.Email{
		Subject: "reset your critiquefi password",
		Body:    html,
		To:      []mail.Recipient{recipient},
	}

	return a.sender.Send(ctx, email)
}
