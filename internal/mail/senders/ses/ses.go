// Package ses provides mail logic using AWS SES.
package ses

import (
	"context"
	"fmt"

	"github.com/andrew-hayworth22/critiquefi-service/internal/mail"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// Config defines the configuration for the SES sender.
type Config struct {
	Region      string
	FromAddress string
	FromName    string
}

// Sender implements the mail.Sender interface using AWS SES.
type Sender struct {
	client      *sesv2.Client
	fromAddress string
	fromName    string
}

// New creates a new SES sender.
func New(ctx context.Context, cfg Config) (*Sender, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing AWS client: %w", err)
	}

	return &Sender{
		client:      sesv2.NewFromConfig(awsCfg),
		fromAddress: cfg.FromAddress,
		fromName:    cfg.FromName,
	}, nil
}

// Send sends an mail using AWS SES.
func (s *Sender) Send(ctx context.Context, email mail.Email) error {
	// Format to addresses
	var toAddresses []string
	for _, r := range email.To {
		toAddresses = append(toAddresses, formatAddress(r.Name, r.Address))
	}

	// Create email input
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(formatAddress(s.fromName, s.fromAddress)),
		Destination: &types.Destination{
			ToAddresses: toAddresses,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(email.Subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(email.Body),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	// Send email
	if _, err := s.client.SendEmail(ctx, input); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}

// formatAddress formats an address as "Name <address>".
func formatAddress(name, address string) string {
	return fmt.Sprintf("%s <%s>", name, address)
}
