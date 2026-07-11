// Package mail provides the abstraction for sending emails
package mail

import (
	"context"
)

// Recipient represents a recipient of an mail
type Recipient struct {
	Name    string
	Address string
}

// Email represents an mail to be sent
type Email struct {
	Subject string
	Body    string
	To      []Recipient
}

// Sender defines the behavior of an mail sender
// Find implementations in the mail/senders directory
type Sender interface {
	Send(ctx context.Context, email Email) error
}
