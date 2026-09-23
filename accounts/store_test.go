package main

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"pgregory.net/rapid"
)

var tokenForm = regexp.MustCompile(`^[A-Z2-7]{26}$`)

// TestTokensProperty checks that tokens are URL-safe, are all different,
// and that different tokens have different hashes.
func TestTokensProperty(t *testing.T) {
	seen := map[string]bool{}
	rapid.Check(t, func(t *rapid.T) {
		tok := newToken()
		if !tokenForm.MatchString(tok) {
			t.Fatalf("token %q is not 26 base32 characters", tok)
		}
		if seen[tok] {
			t.Fatalf("token %q came two times", tok)
		}
		seen[tok] = true
		other := rapid.String().Draw(t, "other")
		if other != tok && hashToken(other) == hashToken(tok) {
			t.Fatal("two tokens have one hash")
		}
	})
}

// TestLoginTokenProperty checks that a sign-in token works one time for
// its address, and that no other string works.
func TestLoginTokenProperty(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		e := newTestEnv(t, nil)
		email := rapid.StringMatching(`[a-z0-9]{1,10}@[a-z]{1,10}\.[a-z]{2,4}`).Draw(rt, "email")
		tok, err := createLoginToken(e.app.db, email, e.clock.now())
		if err != nil || tok == "" {
			rt.Fatalf("createLoginToken: %q, %v", tok, err)
		}
		guess := rapid.String().Filter(func(s string) bool { return s != tok }).Draw(rt, "guess")
		if _, err := consumeLoginToken(e.app.db, guess, e.clock.now()); !errors.Is(err, errNotFound) {
			rt.Fatalf("a wrong token gave %v", err)
		}
		e.advance(time.Duration(rapid.IntRange(0, int(loginTTL/time.Second)-1).Draw(rt, "wait")) * time.Second)
		u, err := consumeLoginToken(e.app.db, tok, e.clock.now())
		if err != nil || u.Email != email {
			rt.Fatalf("consume: %+v, %v", u, err)
		}
		if _, err := consumeLoginToken(e.app.db, tok, e.clock.now()); !errors.Is(err, errNotFound) {
			rt.Fatalf("second use gave %v", err)
		}
	})
}

func TestLoginTokenExpires(t *testing.T) {
	e := newTestEnv(t, nil)
	tok, _ := createLoginToken(e.app.db, "a@example.com", e.clock.now())
	e.advance(loginTTL)
	if _, err := consumeLoginToken(e.app.db, tok, e.clock.now()); !errors.Is(err, errNotFound) {
		t.Fatalf("expired token gave %v", err)
	}
}

func TestLoginTokenSpacing(t *testing.T) {
	e := newTestEnv(t, nil)
	if tok, _ := createLoginToken(e.app.db, "a@example.com", e.clock.now()); tok == "" {
		t.Fatal("first token is empty")
	}
	e.advance(loginSpacing - time.Second)
	if tok, _ := createLoginToken(e.app.db, "a@example.com", e.clock.now()); tok != "" {
		t.Fatal("second token came too soon")
	}
	if tok, _ := createLoginToken(e.app.db, "b@example.com", e.clock.now()); tok == "" {
		t.Fatal("other address got no token")
	}
	e.advance(time.Second)
	if tok, _ := createLoginToken(e.app.db, "a@example.com", e.clock.now()); tok == "" {
		t.Fatal("no token after the spacing")
	}
}

func TestSignInTwiceKeepsOneUser(t *testing.T) {
	e := newTestEnv(t, nil)
	a := e.signUp(t, "a@example.com")
	e.advance(loginSpacing)
	b := e.signUp(t, "a@example.com")
	if a.ID != b.ID || a.FeedToken != b.FeedToken || a.CreatedAt != b.CreatedAt {
		t.Fatalf("second sign-in changed the user: %+v, %+v", a, b)
	}
}

func TestSessionExpires(t *testing.T) {
	e := newTestEnv(t, nil)
	u := e.signUp(t, "a@example.com")
	s, err := createSession(e.app.db, u.ID, e.clock.now())
	if err != nil {
		t.Fatal(err)
	}
	if got, err := sessionUser(e.app.db, s, e.clock.now()); err != nil || got.ID != u.ID {
		t.Fatalf("sessionUser: %+v, %v", got, err)
	}
	e.advance(sessionTTL)
	if _, err := sessionUser(e.app.db, s, e.clock.now()); !errors.Is(err, errNotFound) {
		t.Fatalf("expired session gave %v", err)
	}
}
