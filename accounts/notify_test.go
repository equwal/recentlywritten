package main

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"pgregory.net/rapid"
)

func TestNotifyFirstRunSendsNoOldPosts(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	e.signUp(t, "a@example.com")
	e.advance(time.Minute)
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	if len(e.mail.sent) != 0 {
		t.Fatalf("sent %d messages for posts that were on the site before the first run", len(e.mail.sent))
	}
}

func TestNotifySendsNewPostOneTime(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	if err := e.app.notify(); err != nil { // first run records the old posts
		t.Fatal(err)
	}
	u := e.signUp(t, "a@example.com")
	e.advance(time.Hour)
	newPost := Post{Date: "2026-09-02", Slug: "new-post", Tags: []string{"articles"}, Title: "New post"}
	e.writeIndex(t, append([]Post{newPost}, oldPosts...))

	for range 3 {
		if err := e.app.notify(); err != nil {
			t.Fatal(err)
		}
		e.advance(time.Minute)
	}
	if len(e.mail.sent) != 1 {
		t.Fatalf("sent %d messages, want 1", len(e.mail.sent))
	}
	m := e.mail.sent[0]
	if m.To != "a@example.com" || m.Subject != "New post: New post" {
		t.Errorf("message = %+v", m)
	}
	if !strings.Contains(m.Body, "https://example.com/new-post.html") {
		t.Errorf("body has no link to the post:\n%s", m.Body)
	}
	if strings.Contains(m.Body, "old-article") {
		t.Errorf("body has an old post:\n%s", m.Body)
	}
	wantUnsub := "https://example.com/account/unsubscribe?t=" + u.UnsubToken
	if m.Unsubscribe != wantUnsub || !strings.Contains(m.Body, wantUnsub) {
		t.Errorf("unsubscribe link missing: header %q, body:\n%s", m.Unsubscribe, m.Body)
	}
}

func TestNotifyRetriesAfterSendFailure(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	e.signUp(t, "a@example.com")
	e.advance(time.Hour)
	e.writeIndex(t, append([]Post{{Date: "2026-09-02", Slug: "new-post", Title: "New post"}}, oldPosts...))

	e.mail.fail = true
	if err := e.app.notify(); err == nil {
		t.Fatal("notify returned no error when the send failed")
	}
	e.mail.fail = false
	for range 2 {
		if err := e.app.notify(); err != nil {
			t.Fatal(err)
		}
	}
	if len(e.mail.sent) != 1 {
		t.Fatalf("sent %d messages, want 1", len(e.mail.sent))
	}
}

func TestNotifySkipsPostsFromBeforeSignUp(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	e.advance(time.Hour)
	e.writeIndex(t, append([]Post{{Date: "2026-09-02", Slug: "earlier", Title: "Earlier"}}, oldPosts...))
	if err := e.app.notify(); err != nil { // no users yet
		t.Fatal(err)
	}
	e.advance(time.Hour)
	e.signUp(t, "late@example.com")
	e.advance(time.Hour)
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	if len(e.mail.sent) != 0 {
		t.Fatalf("a new user got a post from before the sign-up: %+v", e.mail.sent)
	}
}

func TestNotifyObeysSettings(t *testing.T) {
	e := newTestEnv(t, oldPosts)
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	u := e.signUp(t, "a@example.com")
	if err := saveSettings(e.app.db, u.ID, true, []string{"projects"}); err != nil {
		t.Fatal(err)
	}
	e.advance(time.Hour)
	e.writeIndex(t, append([]Post{
		{Date: "2026-09-02", Slug: "proj", Tags: []string{"projects"}, Title: "Proj"},
		{Date: "2026-09-02", Slug: "art", Tags: []string{"articles"}, Title: "Art"},
	}, oldPosts...))
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	if len(e.mail.sent) != 1 || strings.Contains(e.mail.sent[0].Body, "proj.html") {
		t.Fatalf("muted tag was sent: %+v", e.mail.sent)
	}

	// Unmute the tag. The muted post must not arrive later.
	if err := saveSettings(e.app.db, u.ID, true, nil); err != nil {
		t.Fatal(err)
	}
	e.advance(time.Hour)
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	if len(e.mail.sent) != 1 {
		t.Fatalf("unmute sent an old post: %+v", e.mail.sent[1:])
	}

	// With email off, nothing goes out.
	if err := saveSettings(e.app.db, u.ID, false, nil); err != nil {
		t.Fatal(err)
	}
	e.advance(time.Hour)
	e.writeIndex(t, append([]Post{{Date: "2026-09-03", Slug: "third", Title: "Third"}}, oldPosts...))
	if err := e.app.notify(); err != nil {
		t.Fatal(err)
	}
	if len(e.mail.sent) != 1 {
		t.Fatalf("email off, but sent: %+v", e.mail.sent[1:])
	}
}

// TestNotifyNoDuplicatesProperty runs random sequences of new posts, new
// users, send failures, and notify runs. At the end no user has a post two
// times, and each user has each post that notify first saw at or after
// the sign-up.
func TestNotifyNoDuplicatesProperty(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		e := newTestEnv(t, oldPosts)
		if err := e.app.notify(); err != nil { // the first run records the old posts
			rt.Fatal(err)
		}
		index := append([]Post(nil), oldPosts...)
		seenAt := map[string]time.Time{}   // slug -> time of the first notify run that saw it
		joinedAt := map[string]time.Time{} // email -> sign-up time
		run := func() {
			for _, p := range index {
				if _, ok := seenAt[p.Slug]; !ok {
					seenAt[p.Slug] = e.clock.now()
				}
			}
			_ = e.app.notify()
		}
		for _, p := range oldPosts {
			seenAt[p.Slug] = time.Time{}
		}

		steps := rapid.SliceOfN(rapid.IntRange(0, 3), 1, 30).Draw(rt, "steps")
		for i, step := range steps {
			e.advance(time.Minute)
			switch step {
			case 0:
				p := Post{Date: "2026-09-02", Slug: fmt.Sprintf("p%d", i), Title: fmt.Sprintf("P%d", i)}
				index = append([]Post{p}, index...)
				e.writeIndex(rt, index)
			case 1:
				email := fmt.Sprintf("u%d@example.com", i)
				e.signUp(rt, email)
				joinedAt[email] = e.clock.now()
			case 2:
				e.mail.fail = !e.mail.fail
			case 3:
				run()
			}
		}
		e.advance(time.Minute)
		e.mail.fail = false
		run()

		got := map[string]map[string]int{} // email -> slug -> count
		for _, m := range e.mail.sent {
			if got[m.To] == nil {
				got[m.To] = map[string]int{}
			}
			for _, line := range strings.Split(m.Body, "\n") {
				if s, ok := strings.CutPrefix(line, "https://example.com/"); ok {
					got[m.To][strings.TrimSuffix(s, ".html")]++
				}
			}
		}
		for email, joined := range joinedAt {
			want := 0
			for _, p := range index {
				if seenAt[p.Slug].Before(joined) {
					continue
				}
				want++
				if n := got[email][p.Slug]; n != 1 {
					rt.Fatalf("%s got %s %d times, want 1", email, p.Slug, n)
				}
			}
			if len(got[email]) != want {
				rt.Fatalf("%s got %v, want %d posts", email, got[email], want)
			}
		}
	})
}
