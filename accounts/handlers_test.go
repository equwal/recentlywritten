package main

import (
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
)

type browser struct {
	t      *testing.T
	client *http.Client
	base   string
}

// newBrowser starts the app on a test server and returns a client with a
// cookie jar.
func newBrowser(t *testing.T, e *testEnv) *browser {
	srv := httptest.NewServer(e.app.routes())
	t.Cleanup(srv.Close)
	e.app.cfg.BaseURL = srv.URL
	jar, _ := cookiejar.New(nil)
	return &browser{t: t, client: &http.Client{Jar: jar}, base: srv.URL}
}

func (b *browser) do(method, path string, form url.Values) (int, string) {
	b.t.Helper()
	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequest(method, b.base+path, body)
	if err != nil {
		b.t.Fatal(err)
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := b.client.Do(req)
	if err != nil {
		b.t.Fatal(err)
	}
	defer resp.Body.Close()
	text, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(text)
}

var linkToken = regexp.MustCompile(`/account/verify\?t=([A-Z2-7]+)`)

// signIn does the sign-in steps of a person in a browser and returns the
// settings page.
func (b *browser) signIn(e *testEnv, email string) string {
	b.t.Helper()
	n := len(e.mail.sent)
	if code, page := b.do("POST", "/account/signin", url.Values{"email": {email}}); code != 200 || !strings.Contains(page, "Check your email") {
		b.t.Fatalf("signin: %d\n%s", code, page)
	}
	if len(e.mail.sent) != n+1 {
		b.t.Fatalf("sent %d sign-in messages, want 1", len(e.mail.sent)-n)
	}
	m := linkToken.FindStringSubmatch(e.mail.sent[n].Body)
	if m == nil {
		b.t.Fatalf("no link in:\n%s", e.mail.sent[n].Body)
	}
	// A GET request, as from a mail scanner, must not use the token.
	if code, page := b.do("GET", "/account/verify?t="+m[1], nil); code != 200 || !strings.Contains(page, `value="`+m[1]+`"`) {
		b.t.Fatalf("verify page: %d\n%s", code, page)
	}
	code, page := b.do("POST", "/account/verify", url.Values{"t": {m[1]}})
	if code != 200 || !strings.Contains(page, "Signed in as "+strings.ToLower(email)) {
		b.t.Fatalf("verify: %d\n%s", code, page)
	}
	return page
}

var feedURL = regexp.MustCompile(`/account/feed/[A-Z2-7]+`)

func TestSignUpAndPrivateFeed(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	b := newBrowser(t, e)

	if code, page := b.do("GET", "/account/", nil); code != 200 || !strings.Contains(page, `action="/account/signin"`) {
		t.Fatalf("signed-out page: %d\n%s", code, page)
	}
	page := b.signIn(e, "Reader@Example.com")
	if e.mail.sent[0].To != "reader@example.com" {
		t.Errorf("sign-in mail went to %q", e.mail.sent[0].To)
	}

	feed := feedURL.FindString(page)
	if feed == "" {
		t.Fatalf("no feed URL on the settings page:\n%s", page)
	}
	// The feed needs no cookie, because a feed reader app has none.
	anon := &browser{t: t, client: http.DefaultClient, base: b.base}
	code, doc := anon.do("GET", feed, nil)
	if code != 200 || wellFormed([]byte(doc)) != nil {
		t.Fatalf("feed: %d\n%s", code, doc)
	}
	for _, want := range []string{"<title>Old article</title>", "old-project.html", "Body of old-article &amp;amp; more."} {
		if !strings.Contains(doc, want) {
			t.Errorf("feed has no %q:\n%s", want, doc)
		}
	}

	// Mute "projects". The feed then leaves out the project post.
	b.do("POST", "/account/settings", url.Values{"email_on": {"1"}, "tag": {"articles"}})
	if _, doc := anon.do("GET", feed, nil); strings.Contains(doc, "old-project.html") || !strings.Contains(doc, "old-article.html") {
		t.Errorf("muted tag still in feed:\n%s", doc)
	}

	// A new feed URL revokes the old one.
	_, page = b.do("POST", "/account/feed/rotate", nil)
	newFeed := feedURL.FindString(page)
	if newFeed == "" || newFeed == feed {
		t.Fatalf("feed URL did not change: %q", newFeed)
	}
	if code, _ := anon.do("GET", feed, nil); code != 404 {
		t.Errorf("old feed URL gave %d, want 404", code)
	}
	if code, _ := anon.do("GET", newFeed, nil); code != 200 {
		t.Errorf("new feed URL gave %d, want 200", code)
	}

	// After sign-out, the settings page asks for an address again.
	if _, page := b.do("POST", "/account/signout", nil); !strings.Contains(page, `action="/account/signin"`) {
		t.Errorf("still signed in after sign-out:\n%s", page)
	}
}

func TestSignInRejectsBadAddress(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	b := newBrowser(t, e)
	for _, bad := range []string{"", "no-at-sign", "Name <a@example.com>", "a@example.com\r\nBcc: x@example.com"} {
		if code, _ := b.do("POST", "/account/signin", url.Values{"email": {bad}}); code != 400 {
			t.Errorf("address %q gave %d, want 400", bad, code)
		}
	}
	if len(e.mail.sent) != 0 {
		t.Errorf("sent mail for a bad address: %+v", e.mail.sent)
	}
}

func TestUsedLinkFails(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	b := newBrowser(t, e)
	b.signIn(e, "a@example.com")
	tok := linkToken.FindStringSubmatch(e.mail.sent[0].Body)[1]
	if code, page := b.do("POST", "/account/verify", url.Values{"t": {tok}}); code != 400 || !strings.Contains(page, "not valid") {
		t.Errorf("second use: %d\n%s", code, page)
	}
}

func TestUnsubscribe(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	b := newBrowser(t, e)
	u := e.signUp(t, "a@example.com")

	// The link in the email body shows a button, and changes nothing.
	if code, page := b.do("GET", "/account/unsubscribe?t="+u.UnsubToken, nil); code != 200 || !strings.Contains(page, "Stop the emails") {
		t.Fatalf("unsubscribe page: %d\n%s", code, page)
	}
	if on := emailOn(t, e, u.ID); !on {
		t.Fatal("GET request turned off email")
	}

	// RFC 8058 one-click: the mail client sends a POST to the header URL.
	code, _ := b.do("POST", "/account/unsubscribe?t="+u.UnsubToken, url.Values{"List-Unsubscribe": {"One-Click"}})
	if code != 200 {
		t.Fatalf("one-click unsubscribe gave %d", code)
	}
	if emailOn(t, e, u.ID) {
		t.Fatal("email is still on after unsubscribe")
	}

	if code, _ := b.do("POST", "/account/unsubscribe", url.Values{"t": {"WRONG"}}); code != 400 {
		t.Errorf("wrong token gave %d, want 400", code)
	}
}

func emailOn(t *testing.T, e *testEnv, id int64) bool {
	t.Helper()
	var on bool
	if err := e.app.db.QueryRow(`SELECT email_on FROM users WHERE id = ?`, id).Scan(&on); err != nil {
		t.Fatal(err)
	}
	return on
}

func TestDeleteAccount(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	b := newBrowser(t, e)
	page := b.signIn(e, "a@example.com")
	feed := feedURL.FindString(page)
	u, err := userByFeedToken(e.app.db, strings.TrimPrefix(feed, "/account/feed/"))
	if err != nil {
		t.Fatal(err)
	}
	if err := saveSettings(e.app.db, u.ID, true, []string{"projects"}); err != nil {
		t.Fatal(err)
	}
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}

	if code, page := b.do("POST", "/account/delete", nil); code != 200 || !strings.Contains(page, "Account deleted") {
		t.Fatalf("delete: %d\n%s", code, page)
	}
	for _, table := range []string{"users", "sessions", "muted_tags", "deliveries"} {
		var n int
		if err := e.app.db.QueryRow(`SELECT count(*) FROM ` + table).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n != 0 {
			t.Errorf("%s has %d rows after delete", table, n)
		}
	}
	if code, _ := b.do("GET", feed, nil); code != 404 {
		t.Errorf("feed after delete gave %d, want 404", code)
	}
	if _, err := userByFeedToken(e.app.db, u.FeedToken); !errors.Is(err, errNotFound) {
		t.Errorf("user still there: %v", err)
	}
}

func TestStateChangesNeedSession(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	u := e.signUp(t, "a@example.com")
	b := newBrowser(t, e)
	// No cookie: each POST goes to the sign-in page and changes nothing.
	for _, path := range []string{"/account/delete", "/account/settings", "/account/feed/rotate"} {
		if _, page := b.do("POST", path, url.Values{}); !strings.Contains(page, `action="/account/signin"`) {
			t.Errorf("%s without a session:\n%s", path, page)
		}
	}
	after, err := userByFeedToken(e.app.db, u.FeedToken)
	if err != nil || !after.EmailOn {
		t.Errorf("user changed: %+v, %v", after, err)
	}
}
