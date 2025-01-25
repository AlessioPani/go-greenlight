package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// Set up embedding templates into the executable.

//go:embed "templates"
var templateFS embed.FS

// Mailer struct which contains a mail.Dialer instance (used to connect to a
// SMTP server) and the sender information.
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// New initialize a new Mailer instance.
func New(host string, port int, username string, password string, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send send an email using a specified template file.
func (m Mailer) Send(recipient string, templateFile string, data any) error {
	// Parse a template from our embedded FS.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the named template "subject" and stores it into a buffer.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Execute the named template "plainBody" and stores it into a buffer.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Execute the named template "htmlBody" and stores it into a buffer.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Prepare a message.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	// Connect to the configured SMTP server, send the message and close the connection.
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}

	return nil
}
