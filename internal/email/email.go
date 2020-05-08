package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	mailer "github.com/jordan-wright/email"
)

// Sender is the email client that allows to send the emails
type Sender struct {
	auth     smtp.Auth
	hostname string
	port     int
	from     string
}

// NewSender creates a new email client
func NewSender(username, password, hostname, from string, port int) *Sender {
	s := Sender{
		auth:     smtp.PlainAuth("", username, password, hostname),
		hostname: hostname,
		port:     port,
		from:     from,
	}
	return &s
}

//SendInvitation will send an invitation email to a specific email
func (s *Sender) SendInvitation(to string, url string) error {
	var data struct {
		URL string
	}
	data.URL = url
	buf := new(bytes.Buffer)
	if err := parseTemplate("message-invitation.html", data, buf); err == nil {
		return s.send(to, "WG Concierge Invitation", buf.Bytes())
	}
	return nil
}

func (s *Sender) fullHostname() string {
	return fmt.Sprintf("%s:%d", s.hostname, s.port)
}

func (s *Sender) send(to, subject string, body []byte) error {
	e := mailer.NewEmail()
	e.To = []string{to}
	e.From = s.from
	e.Subject = subject
	e.HTML = body
	return e.Send(s.fullHostname(), s.auth)
}

func parseTemplate(templateName string, data interface{}, buf *bytes.Buffer) error {
	t, err := template.New(templateName).Parse(messageInvitation)
	if err != nil {
		return err
	}
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	return nil
}

const messageInvitation = `
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html>
	<head/>
	<body>
		<p>
    		Activate your WireGuard Client on {{ .URL }}
		</p>
	</body>
</html>
`
