package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func NewMailer(host string, port int, username, pasword, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, pasword)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}
	// Execute the named template "subject", passing in the dynamic data and storing the
	// result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())
	// Call the DialAndSend() method on the dialer, passing in the message to send. This
	// opens a connection to the SMTP server, sends the message, then closes the
	// connection. If there is a timeout, it will return a "dial tcp: i/o timeout"
	// error.
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
