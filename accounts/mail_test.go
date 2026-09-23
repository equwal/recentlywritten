package main

import (
	"bytes"
	"io"
	"mime"
	"mime/quotedprintable"
	"net/mail"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"pgregory.net/rapid"
)

var testFrom = &mail.Address{Name: "Recently Written", Address: "posts@example.com"}

// TestFormatMessageProperty checks that any subject and body give a
// message that parses, has exactly the expected headers, and decodes back
// to the same subject and body.
func TestFormatMessageProperty(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		valid := rapid.String().Filter(utf8.ValidString)
		m := message{
			To:          "reader@example.com",
			Subject:     valid.Draw(t, "subject"),
			Body:        valid.Draw(t, "body"),
			Unsubscribe: rapid.SampledFrom([]string{"", "https://example.com/account/unsubscribe?t=ABC"}).Draw(t, "unsub"),
		}
		data, err := formatMessage(testFrom, m, time.Unix(0, 0))
		if err != nil {
			t.Fatal(err)
		}
		msg, err := mail.ReadMessage(bytes.NewReader(data))
		if err != nil {
			t.Fatalf("message does not parse: %v\n%q", err, data)
		}
		wantHeaders := 8
		if m.Unsubscribe != "" {
			wantHeaders += 2
		}
		if len(msg.Header) != wantHeaders {
			t.Fatalf("got headers %v, want %d (a line break in the subject must not add one)", msg.Header, wantHeaders)
		}
		subject, err := new(mime.WordDecoder).DecodeHeader(msg.Header.Get("Subject"))
		if err != nil {
			t.Fatal(err)
		}
		wantSubject := strings.NewReplacer("\r", " ", "\n", " ").Replace(m.Subject)
		// net/mail removes spaces at the two ends of a plain header value,
		// but not in an encoded word.
		if strings.Trim(subject, " \t") != strings.Trim(wantSubject, " \t") {
			t.Fatalf("subject: got %q, want %q", subject, wantSubject)
		}
		body, err := io.ReadAll(quotedprintable.NewReader(msg.Body))
		if err != nil {
			t.Fatal(err)
		}
		lf := strings.NewReplacer("\r\n", "\n", "\r", "\n")
		if got := lf.Replace(string(body)); got != lf.Replace(m.Body) {
			t.Fatalf("body: got %q, want %q", got, m.Body)
		}
		if got := msg.Header.Get("List-Unsubscribe"); m.Unsubscribe != "" && got != "<"+m.Unsubscribe+">" {
			t.Fatalf("List-Unsubscribe = %q", got)
		}
	})
}

func TestFormatMessageRejectsBadAddress(t *testing.T) {
	if _, err := formatMessage(testFrom, message{To: "not an address"}, time.Now()); err == nil {
		t.Fatal("no error for a bad address")
	}
}

func TestLoadConfigNamesMissingVariables(t *testing.T) {
	env := map[string]string{"RW_BASE_URL": "https://example.com/", "SMTP_PORT": "2525"}
	_, err := loadConfig(func(k string) string { return env[k] }, false)
	if err == nil || err.Error() != "set these environment variables: RW_SITE_DIR, RW_DB_PATH, RW_MAIL_FROM, SMTP_HOST" {
		t.Fatalf("err = %v", err)
	}

	env["RW_SITE_DIR"], env["RW_DB_PATH"] = "/srv/site", "/var/lib/rw/rw.db"
	env["RW_MAIL_FROM"], env["SMTP_HOST"] = "Recently Written <posts@example.com>", "smtp.example.com"
	c, err := loadConfig(func(k string) string { return env[k] }, false)
	if err != nil {
		t.Fatal(err)
	}
	if c.BaseURL != "https://example.com" || c.SMTPAddr != "smtp.example.com:2525" ||
		c.MailFrom.Address != "posts@example.com" || c.Listen != "127.0.0.1:8081" {
		t.Fatalf("config = %+v", c)
	}

	delete(env, "SMTP_HOST")
	if _, err := loadConfig(func(k string) string { return env[k] }, true); err != nil {
		t.Fatalf("dev mode needs no SMTP: %v", err)
	}
}
