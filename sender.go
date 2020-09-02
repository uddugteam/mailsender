package mailsender

import (
	"encoding/base64"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

// Smp represent SMTP server mail sender
type Smtp struct {
	Addr string
}

// NewSmtp create new SMTP mail sender
func NewSmtp(addr string) *Smtp {
	return &Smtp{Addr: addr}
}

// Send send mail with SMTP server
func (s *Smtp) Send(from, to, subject, body string) error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(s.Addr)
	if err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}

	msg := composeMimeMail(to, from, subject, body)

	if _, err := wc.Write(msg); err != nil {
		return err
	}

	if err = wc.Close(); err != nil {
		return err
	}

	if err = c.Quit(); err != nil {
		return err
	}

	return nil
}

// compose mail body
func composeMimeMail(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = formatEmailAddress(from)
	header["To"] = formatEmailAddress(to)
	header["Subject"] = encodeRFC2047(subject)
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
	header["Date"] = time.Now().String()

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	base64Text := make([]byte, base64.StdEncoding.EncodedLen(len(body)))
	base64.StdEncoding.Encode(base64Text, []byte(body))

	message += "\r\n" + string(base64Text)

	return []byte(message)
}

// Never fails, tries to format the address if possible
func formatEmailAddress(addr string) string {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return addr
	}
	return e.String()
}

// encode string to rfc2047
func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Address: str}
	return strings.TrimSuffix(strings.Trim(addr.String(), " <>"), "@")
}
