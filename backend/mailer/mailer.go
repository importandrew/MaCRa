package mailer

import (
	"fmt"
	"net"
	"net/smtp"
	"net/url"
	"strings"

	"github.com/importandrew/MaCRa/models"
)

const (
	smtpHost = "127.0.0.1"
	smtpPort = "25"
)

type Mailer struct {
	from string
}

func New(from string) *Mailer {
	return &Mailer{from: from}
}

func (m *Mailer) SendVerification(to, origin, token string) error {
	link := origin + "/verify-email?token=" + url.QueryEscape(token)

	body := fmt.Sprintf("Verify your MaCRa account by opening this link:\n\n%s\n\nThis link expires in %s.\n",
		link, models.EmailVerifyTTL.String(),
	)

	return m.send(to, "Verify your MaCRa email", body)
}

func (m *Mailer) SendPasswordReset(to, origin, token string) error {
	link := origin + "/reset-password?token=" + url.QueryEscape(token)

	body := fmt.Sprintf("Reset your MaCRa password by opening this link:\n\n%s\n\nThis link expires in %s.\nIf you did not request a reset, ignore this email.\n",
		link, models.PasswordResetTTL.String(),
	)

	return m.send(to, "Reset your MaCRa password", body)
}

// SendAccountExists answers a sign-up with an address that already has an
// account: it says so and offers a password reset, instead of the sign-up
// telling a stranger the address is registered.
func (m *Mailer) SendAccountExists(to, origin, token string) error {
	link := origin + "/reset-password?token=" + url.QueryEscape(token)

	body := fmt.Sprintf("Someone tried to sign up for MaCRa with this email, but it already has an account.\n\n"+
		"If it was you, sign in with your password, or set a new one by opening this link:\n\n%s\n\n"+
		"This link expires in %s.\nIf it wasn't you, ignore this email: nothing has changed.\n",
		link, models.PasswordResetTTL.String(),
	)

	return m.send(to, "You already have a MaCRa account", body)
}

func (m *Mailer) send(to, subject, body string) error {
	msg := strings.Join([]string{
		"From: " + m.from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	addr := net.JoinHostPort(smtpHost, smtpPort)
	return smtp.SendMail(addr, nil, m.from, []string{to}, []byte(msg))
}
