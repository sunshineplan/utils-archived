package mail

import (
	"bytes"
	"os"
	"path/filepath"
	"time"

	"github.com/go-mail/mail"
)

// Dialer is a dialer to an SMTP server.
type Dialer struct {
	Host     string
	Port     int
	Account  string
	Password string
	Timeout  time.Duration
}

// Message represents an email.
type Message struct {
	To, Cc, Bcc []string
	Subject     string
	Body        string
	Attachments []*Attachment
}

// Attachment represents an attachment.
type Attachment struct {
	Filename string
	Path     string
	Bytes    []byte
}

// Send sends the given messages.
func (d *Dialer) Send(msg ...*Message) error {
	for _, m := range msg {
		message := mail.NewMessage()
		message.SetHeader("From", d.Account)
		message.SetHeader("To", m.To...)
		message.SetHeader("Cc", m.Cc...)
		message.SetHeader("Bcc", m.Bcc...)
		message.SetHeader("Subject", m.Subject)
		message.SetBody("text/plain", m.Body)
		for _, a := range m.Attachments {
			if a.Bytes != nil {
				if a.Filename == "" {
					a.Filename = "attachment"
				}
				message.AttachReader(a.Filename, bytes.NewBuffer(a.Bytes))
				continue
			}
			if a.Filename == "" {
				a.Filename = filepath.Base(a.Path)
			}
			f, err := os.Open(a.Path)
			if err != nil {
				return err
			}
			message.AttachReader(a.Filename, f)
		}

		dialer := mail.NewDialer(d.Host, d.Port, d.Account, d.Password)
		if d.Timeout != 0 {
			dialer.Timeout = d.Timeout
		} else {
			dialer.Timeout = 3 * time.Minute
		}

		if err := dialer.DialAndSend(message); err != nil {
			return err
		}
	}

	return nil
}
