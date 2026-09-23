package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// captureMailer keeps each message. If fail is true, Send returns an error
// and keeps nothing.
type captureMailer struct {
	mu   sync.Mutex
	sent []message
	fail bool
}

func (c *captureMailer) Send(m message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.fail {
		return errors.New("smtp is down")
	}
	c.sent = append(c.sent, m)
	return nil
}

// clock is a time source that a test moves forward.
type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

type testEnv struct {
	app   *app
	mail  *captureMailer
	clock *clock
}

// fataler is the part of testing.TB that the helpers use. *rapid.T has
// these methods too.
type fataler interface {
	Helper()
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

// newTestEnv makes an app with a :memory: database and a site directory
// that has the posts in index.
func newTestEnv(t *testing.T, index []Post) *testEnv {
	t.Helper()
	db, err := openDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	site := t.TempDir()
	if err := os.MkdirAll(filepath.Join(site, "postdata"), 0o755); err != nil {
		t.Fatal(err)
	}
	e := &testEnv{
		mail:  &captureMailer{},
		clock: &clock{t: time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)},
	}
	e.app = &app{
		cfg:  config{BaseURL: "https://example.com", SiteDir: site},
		db:   db,
		mail: e.mail,
		log:  log.New(io.Discard, "", 0),
	}
	e.app.now = e.clock.now
	e.writeIndex(t, index)
	return e
}

// writeIndex writes postdata/index.tsv and a body for each post, as
// build.sh does.
func (e *testEnv) writeIndex(t fataler, posts []Post) {
	t.Helper()
	var b strings.Builder
	for _, p := range posts {
		b.WriteString(p.Date + "\t" + p.Slug + "\t" + strings.Join(p.Tags, " ") + "\t" + p.Title + "\n")
		body := "<p>Body of " + p.Slug + " &amp; more.</p>\n"
		if err := os.WriteFile(filepath.Join(e.app.cfg.SiteDir, "postdata", p.Slug+".html"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(e.app.cfg.SiteDir, "postdata", "index.tsv"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
}

// signUp makes a verified user for email and returns it.
func (e *testEnv) signUp(t fataler, email string) user {
	t.Helper()
	tok, err := createLoginToken(e.app.db, email, e.clock.now())
	if err != nil || tok == "" {
		t.Fatalf("createLoginToken: %q, %v", tok, err)
	}
	u, err := consumeLoginToken(e.app.db, tok, e.clock.now())
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func (e *testEnv) advance(d time.Duration) { e.clock.t = e.clock.t.Add(d) }

var oldPosts = []Post{
	{Date: "2025-09-11", Slug: "old-article", Tags: []string{"articles"}, Title: "Old article"},
	{Date: "2025-09-06", Slug: "old-project", Tags: []string{"projects"}, Title: "Old project"},
}
