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
