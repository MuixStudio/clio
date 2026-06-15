// Package courier sends transactional messages (verification codes, recovery links, …).
//
// The package defines one Courier interface so callers never depend on a concrete transport:
//   - ConsoleCourier — prints to stdout; zero config, ideal for development.
//   - SMTPCourier    — sends real email via net/smtp.
//
// The Driver holds the active Courier and injects it into any strategy that needs it.
package courier

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
)

// Courier is the single interface all strategies use to send messages.
// Swap the implementation in the Driver without touching strategy code.
type Courier interface {
	// SendVerificationCode delivers a one-time code to the given address.
	SendVerificationCode(ctx context.Context, to, code string) error
}

// =========================================================================
// ConsoleCourier — development / testing
// =========================================================================

// ConsoleCourier logs every message to stdout instead of sending email.
// Use it locally so you can see codes without an SMTP server.
type ConsoleCourier struct{}

func NewConsoleCourier() *ConsoleCourier { return &ConsoleCourier{} }

func (c *ConsoleCourier) SendVerificationCode(_ context.Context, to, code string) error {
	log.Printf("[courier] verification code for %s → %s (valid 10 min)", to, code)
	return nil
}

// =========================================================================
// SMTPCourier — production
// =========================================================================

// SMTPConfig holds credentials for an SMTP relay.
type SMTPConfig struct {
	Host     string // e.g. "smtp.sendgrid.net"
	Port     int    // e.g. 587
	Username string
	Password string
	From     string // e.g. "noreply@example.com"
}

// SMTPCourier sends email via net/smtp (works with any SMTP relay).
type SMTPCourier struct{ cfg SMTPConfig }

func NewSMTPCourier(cfg SMTPConfig) *SMTPCourier { return &SMTPCourier{cfg: cfg} }

func (c *SMTPCourier) SendVerificationCode(_ context.Context, to, code string) error {
	addr := fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port)
	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: Your verification code\r\n\r\n"+
			"Your verification code is: %s\r\n\r\nThis code expires in 10 minutes.\r\n",
		c.cfg.From, to, code,
	)
	return smtp.SendMail(addr, auth, c.cfg.From, []string{to}, []byte(msg))
}
