package main

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

// message is one plain text email to one address.
type message struct {
	To      string
	Subject string
	Body    string
	// Unsubscribe is a URL for the List-Unsubscribe header (RFC 2369) and
	// for one-click unsubscribe (RFC 8058). It is empty for sign-in email.
	Unsubscribe string
}

type mailer interface {
	Send(m message) error
}

// smtpMailer sends through any SMTP server: a local MTA, or the SMTP
// endpoint of a hosted provider (Resend, Postmark, Amazon SES).
type smtpMailer struct {
	Addr     string // host:port
	Username string // empty for a server that needs no sign-in
	Password string
	From     *mail.Address
}

// Send uses STARTTLS when the server offers it. net/smtp PlainAuth does
// not send the password on a connection without TLS, except to localhost.
func (s smtpMailer) Send(m message) error {
	var auth smtp.Auth
	if s.Username != "" {
		host, _, _ := net.SplitHostPort(s.Addr)
		auth = smtp.PlainAuth("", s.Username, s.Password, host)
	}
	data, err := formatMessage(s.From, m, time.Now())
	if err != nil {
		return err
	}
	return smtp.SendMail(s.Addr, auth, s.From.Address, []string{m.To}, data)
}

// logMailer writes each message to w. Use it only for development: the
// messages contain sign-in links.
type logMailer struct{ W io.Writer }

func (l logMailer) Send(m message) error {
	_, err := fmt.Fprintf(l.W, "--- mail to %s\nSubject: %s\n\n%s\n---\n", m.To, m.Subject, m.Body)
	return err
}

// formatMessage makes an RFC 5322 message. The subject is encoded as an
// RFC 2047 word, and the body as quoted-printable, so that non-ASCII text
// and long lines arrive intact.
func formatMessage(from *mail.Address, m message, now time.Time) ([]byte, error) {
	to, err := mail.ParseAddress(m.To)
	if err != nil {
		return nil, err
	}
	// A line break in a header value would start a new header.
	subject := strings.NewReplacer("\r", " ", "\n", " ").Replace(m.Subject)
	domain := from.Address[strings.LastIndex(from.Address, "@")+1:]

	var b bytes.Buffer
	h := func(k, v string) { fmt.Fprintf(&b, "%s: %s\r\n", k, v) }
	h("From", from.String())
	h("To", to.String())
	h("Subject", mime.QEncoding.Encode("utf-8", subject))
	h("Date", now.Format(time.RFC1123Z))
	h("Message-ID", "<"+strings.ToLower(newToken())+"@"+domain+">")
	h("MIME-Version", "1.0")
	h("Content-Type", "text/plain; charset=utf-8")
	h("Content-Transfer-Encoding", "quoted-printable")
	if m.Unsubscribe != "" {
		h("List-Unsubscribe", "<"+m.Unsubscribe+">")
		h("List-Unsubscribe-Post", "List-Unsubscribe=One-Click")
	}
	b.WriteString("\r\n")
	qp := quotedprintable.NewWriter(&b)
	// Email lines end with CRLF. A CR or LF alone becomes CRLF.
	crlf := strings.NewReplacer("\r\n", "\r\n", "\r", "\r\n", "\n", "\r\n")
	if _, err := qp.Write([]byte(crlf.Replace(m.Body))); err != nil {
		return nil, err
	}
	if err := qp.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
