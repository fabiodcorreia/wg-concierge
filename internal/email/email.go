package email

import (
	"bytes"
	"html/template"

	mail "github.com/xhit/go-simple-mail"
)

// Sender is the email client that allows to send the emails
type Sender struct {
	server *mail.SMTPServer
	from   string
}

// NewSender creates a new email client
func NewSender(username, password, hostname, from string, port int) *Sender {
	server := mail.NewSMTPClient()
	server.Host = hostname
	server.Port = port
	server.Username = username
	server.Password = password
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	s := Sender{
		server: server,
		from:   from,
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
		return s.send(to, "WG Concierge Invitation", buf.String())
	}
	return nil
}

func (s *Sender) send(to, subject string, body string) error {
	c, err := s.server.Connect()
	if err != nil {
		return err
	}
	defer c.Close()
	message := mail.NewMSG()
	message.SetFrom(s.from).AddTo(to).SetSubject(subject)
	message.SetBody(mail.TextHTML, body)
	return message.Send(c)
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
