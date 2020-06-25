package mail

import (
	"gopkg.in/gomail.v2"
)

// Setting contains mail setting
type Setting struct {
	From           string
	To             []string
	Password       string
	SMTPServer     string
	SMTPServerPort int
}

// Attachment contains attachment file path and filename when sending mail
type Attachment struct {
	FilePath string
	Filename string
}

// SendMail according mail setting, subject, body and attachment
func SendMail(s *Setting, subject string, body string, attachment *Attachment) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", s.To...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	if attachment != nil {
		if attachment.Filename != "" {
			m.Attach(attachment.FilePath, gomail.Rename(attachment.Filename))
		} else {
			m.Attach(attachment.FilePath)
		}
	}

	d := gomail.NewDialer(s.SMTPServer, s.SMTPServerPort, s.From, s.Password)

	err := d.DialAndSend(m)
	return err
}
