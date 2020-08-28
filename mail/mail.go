package mail

import (
	"io"

	"github.com/go-mail/mail"
)

// Setting contains mail setting
type Setting struct {
	From           string
	To, Cc, Bcc    []string
	Password       string
	SMTPServer     string
	SMTPServerPort int
}

// Attachment contains attachment file path and filename when sending mail
type Attachment struct {
	FilePath string
	Filename string
	Reader   io.Reader
}

// Send mail according setting, subject, body and attachment
func (s *Setting) Send(subject string, body string, attachments ...*Attachment) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", s.To...)
	m.SetHeader("Cc", s.Cc...)
	m.SetHeader("Bcc", s.Bcc...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	for _, attachment := range attachments {
		if attachment.Reader != nil {
			m.AttachReader(attachment.Filename, attachment.Reader)
		} else if attachment.Filename != "" {
			m.Attach(attachment.FilePath, mail.Rename(attachment.Filename))
		} else {
			m.Attach(attachment.FilePath)
		}
	}

	d := mail.NewDialer(s.SMTPServer, s.SMTPServerPort, s.From, s.Password)

	err := d.DialAndSend(m)
	return err
}
