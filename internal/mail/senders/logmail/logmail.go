// Package logmail provides mail logic using a logger.
package logmail

import (
	"context"
	"log/slog"

	"github.com/andrew-hayworth22/critiquefi-service/internal/mail"
)

type Sender struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *Sender {
	return &Sender{logger: logger}
}

func (s *Sender) Send(ctx context.Context, email mail.Email) error {
	s.logger.Info("sending email", "subject", email.Subject, "body", email.Body, "recipients", email.To)
	return nil
}
