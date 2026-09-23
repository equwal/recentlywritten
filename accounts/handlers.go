package main

import (
	"bytes"
	"database/sql"
	_ "embed"
	"errors"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

//go:embed templates.html
var templateText string

var pages = template.Must(template.New("").Parse(templateText))

type app struct {
	cfg  config
	db   *sql.DB
	mail mailer
	now  func() time.Time
	log  *log.Logger
}

const sessionCookie = "rw_session"

// routes gives the handler for all paths under /account/. The web server
// in front sends only these paths to this service.
//
// Each change of state uses POST. A GET request from a link in an email
// only shows a page with a button, because mail scanners open the links
// in email and must not use a sign-in link or unsubscribe a user. The
// session cookie is SameSite=Lax, so a form on another site cannot send a
// POST request with the cookie.
func (a *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /account/{$}", a.home)
	mux.HandleFunc("POST /account/signin", a.signin)
	mux.HandleFunc("GET /account/verify", a.verifyPage)
	mux.HandleFunc("POST /account/verify", a.verify)
	mux.HandleFunc("POST /account/settings", a.withUser(a.settings))
	mux.HandleFunc("POST /account/feed/rotate", a.withUser(a.rotateFeed))
	mux.HandleFunc("POST /account/signout", a.signout)
	mux.HandleFunc("POST /account/delete", a.withUser(a.deleteAccount))
	mux.HandleFunc("GET /account/unsubscribe", a.unsubscribePage)
	mux.HandleFunc("POST /account/unsubscribe", a.unsubscribe)
	mux.HandleFunc("GET /account/feed/{token}", a.feed)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The URLs of some pages contain tokens. These headers keep the
		// tokens out of caches and out of the Referer header.
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
		mux.ServeHTTP(w, r)
	})
}

func (a *app) render(w http.ResponseWriter, status int, name string, data any) {
	var b bytes.Buffer
	if err := pages.ExecuteTemplate(&b, name, data); err != nil {
		a.fail(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	b.WriteTo(w)
}

func (a *app) fail(w http.ResponseWriter, err error) {
	a.log.Printf("error: %v", err)
	http.Error(w, "Internal error", http.StatusInternalServerError)
}

func (a *app) currentUser(r *http.Request) (user, error) {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return user{}, errNotFound
	}
	return sessionUser(a.db, c.Value, a.now())
}

func (a *app) withUser(h func(http.ResponseWriter, *http.Request, user)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := a.currentUser(r)
		if errors.Is(err, errNotFound) {
			http.Redirect(w, r, "/account/", http.StatusSeeOther)
			return
		}
		if err != nil {
			a.fail(w, err)
			return
		}
		h(w, r, u)
	}
}

func (a *app) setSessionCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    value,
		Path:     "/account/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   strings.HasPrefix(a.cfg.BaseURL, "https://"),
		SameSite: http.SameSiteLaxMode,
	})
}

type tagChoice struct {
	Name string
	On   bool
}

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	u, err := a.currentUser(r)
	if errors.Is(err, errNotFound) {
		a.render(w, http.StatusOK, "signin", map[string]string{})
		return
	}
	if err != nil {
		a.fail(w, err)
		return
	}
	muted, err := mutedTags(a.db, u.ID)
	if err != nil {
		a.fail(w, err)
		return
	}
	posts, err := loadPosts(a.cfg.SiteDir)
	if err != nil {
		a.fail(w, err)
		return
	}
	var tags []tagChoice
	for _, t := range allTags(posts) {
		tags = append(tags, tagChoice{Name: t, On: !muted[t]})
	}
	a.render(w, http.StatusOK, "settings", map[string]any{
		"Email":   u.Email,
		"EmailOn": u.EmailOn,
		"Tags":    tags,
		"FeedURL": a.cfg.BaseURL + "/account/feed/" + u.FeedToken,
	})
}

// normalizeEmail returns the address in lower case, or an error if s is
// not one plain address such as "name@example.com".
func normalizeEmail(s string) (string, error) {
	s = strings.TrimSpace(s)
	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Name != "" || addr.Address != s || len(s) > 254 {
		return "", errors.New("not a valid address")
	}
	return strings.ToLower(s), nil
}

func (a *app) signin(w http.ResponseWriter, r *http.Request) {
	email, err := normalizeEmail(r.PostFormValue("email"))
	if err != nil {
		a.render(w, http.StatusBadRequest, "signin", map[string]string{"Error": "Type a valid email address."})
		return
	}
	token, err := createLoginToken(a.db, email, a.now())
	if err != nil {
		a.fail(w, err)
		return
	}
	// The page is the same when createLoginToken sends nothing, so the
	// page does not tell which addresses have an account.
	if token != "" {
		link := a.cfg.BaseURL + "/account/verify?t=" + token
		err = a.mail.Send(message{
			To:      email,
			Subject: "Sign in to Recently Written",
			Body: "Open this link to sign in to Recently Written:\n\n" + link +
				"\n\nThe link works one time and expires in 15 minutes.\n" +
				"If you did not ask for this email, ignore it.\n",
		})
		if err != nil {
			a.fail(w, err)
			return
		}
	}
	a.render(w, http.StatusOK, "sent", nil)
}

func (a *app) verifyPage(w http.ResponseWriter, r *http.Request) {
	a.render(w, http.StatusOK, "verify", r.URL.Query().Get("t"))
}

func (a *app) verify(w http.ResponseWriter, r *http.Request) {
	u, err := consumeLoginToken(a.db, r.PostFormValue("t"), a.now())
	if errors.Is(err, errNotFound) {
		a.render(w, http.StatusBadRequest, "badlink", nil)
		return
	}
	if err != nil {
		a.fail(w, err)
		return
	}
	token, err := createSession(a.db, u.ID, a.now())
	if err != nil {
		a.fail(w, err)
		return
	}
	a.setSessionCookie(w, token, int(sessionTTL.Seconds()))
	http.Redirect(w, r, "/account/", http.StatusSeeOther)
}

func (a *app) settings(w http.ResponseWriter, r *http.Request, u user) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	posts, err := loadPosts(a.cfg.SiteDir)
	if err != nil {
		a.fail(w, err)
		return
	}
	checked := map[string]bool{}
	for _, t := range r.PostForm["tag"] {
		checked[t] = true
	}
	// A tag is muted if the form shows it and the user cleared its box.
	var muted []string
	for _, t := range allTags(posts) {
		if !checked[t] {
			muted = append(muted, t)
		}
	}
	if err := saveSettings(a.db, u.ID, r.PostForm.Get("email_on") == "1", muted); err != nil {
		a.fail(w, err)
		return
	}
	http.Redirect(w, r, "/account/", http.StatusSeeOther)
}

func (a *app) rotateFeed(w http.ResponseWriter, r *http.Request, u user) {
	if err := rotateFeedToken(a.db, u.ID); err != nil {
		a.fail(w, err)
		return
	}
	http.Redirect(w, r, "/account/", http.StatusSeeOther)
}

func (a *app) signout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(sessionCookie); err == nil {
		if err := deleteSession(a.db, c.Value); err != nil {
			a.fail(w, err)
			return
		}
	}
	a.setSessionCookie(w, "", -1)
	http.Redirect(w, r, "/account/", http.StatusSeeOther)
}

func (a *app) deleteAccount(w http.ResponseWriter, r *http.Request, u user) {
	if err := deleteUser(a.db, u.ID); err != nil {
		a.fail(w, err)
		return
	}
	a.setSessionCookie(w, "", -1)
	a.render(w, http.StatusOK, "deleted", nil)
}

func (a *app) unsubscribePage(w http.ResponseWriter, r *http.Request) {
	a.render(w, http.StatusOK, "unsubscribe", r.URL.Query().Get("t"))
}

// unsubscribe takes the token from the form, or from the URL. A mail
// client that does RFC 8058 one-click unsubscribe sends a POST request to
// the URL in the List-Unsubscribe header, with the token in the URL.
func (a *app) unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("t")
	if token == "" {
		token = r.URL.Query().Get("t")
	}
	err := unsubscribe(a.db, token)
	if errors.Is(err, errNotFound) {
		a.render(w, http.StatusBadRequest, "badlink", nil)
		return
	}
	if err != nil {
		a.fail(w, err)
		return
	}
	a.render(w, http.StatusOK, "unsubscribed", nil)
}

func (a *app) feed(w http.ResponseWriter, r *http.Request) {
	u, err := userByFeedToken(a.db, r.PathValue("token"))
	if errors.Is(err, errNotFound) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		a.fail(w, err)
		return
	}
	muted, err := mutedTags(a.db, u.ID)
	if err != nil {
		a.fail(w, err)
		return
	}
	posts, err := loadPosts(a.cfg.SiteDir)
	if err != nil {
		a.fail(w, err)
		return
	}
	ch := buildFeed(a.cfg.BaseURL, posts, muted, func(slug string) string {
		return loadBody(a.cfg.SiteDir, slug)
	})
	var b bytes.Buffer
	if err := writeFeed(&b, ch); err != nil {
		a.fail(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	b.WriteTo(w)
}
