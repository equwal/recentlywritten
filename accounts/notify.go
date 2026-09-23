package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// notify sends one email to each subscriber who has new posts.
//
// A post is new for a user if the service first saw it in the index at or
// after the time the user signed up. On the first run the posts table is
// empty, and notify records every post with first_seen 0, so that the
// posts that were on the site before the service started go to nobody.
//
// notify writes a row in deliveries for each new (user, post) before it
// sends the email, and sends only the posts whose rows it wrote. It also
// writes the row when the user does not want the post (email off, or all
// tags of the post muted), so a later change of settings does not send
// old posts. A second run, or a second process at the same time, never
// sends a post two times to one user. If the send fails, notify removes the rows, and the next
// run tries again. If the process stops between the write and the send,
// the user does not get that email. This design accepts a lost email to
// prevent a duplicate email.
func (a *app) notify() error {
	posts, err := loadPosts(a.cfg.SiteDir)
	if err != nil {
		return err
	}
	firstSeen, err := recordPosts(a.db, posts, a.now().Unix())
	if err != nil {
		return err
	}

	rows, err := a.db.Query(`SELECT id, email, created_at, email_on, unsub_token FROM users`)
	if err != nil {
		return err
	}
	var users []user
	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.EmailOn, &u.UnsubToken); err != nil {
			rows.Close()
			return err
		}
		users = append(users, u)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	var errs []error
	for _, u := range users {
		if err := a.notifyUser(u, posts, firstSeen); err != nil {
			errs = append(errs, fmt.Errorf("user %d: %w", u.ID, err))
		}
	}
	return errors.Join(errs...)
}

// recordPosts adds the posts that the posts table does not have, and
// returns first_seen for each post in posts.
func recordPosts(db *sql.DB, posts []Post, now int64) (map[string]int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var known int
	if err := tx.QueryRow(`SELECT count(*) FROM posts`).Scan(&known); err != nil {
		return nil, err
	}
	seen := now
	if known == 0 {
		seen = 0
	}
	firstSeen := map[string]int64{}
	for _, p := range posts {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO posts (slug, first_seen) VALUES (?, ?)`, p.Slug, seen); err != nil {
			return nil, err
		}
		var t int64
		if err := tx.QueryRow(`SELECT first_seen FROM posts WHERE slug = ?`, p.Slug).Scan(&t); err != nil {
			return nil, err
		}
		firstSeen[p.Slug] = t
	}
	return firstSeen, tx.Commit()
}

func (a *app) notifyUser(u user, posts []Post, firstSeen map[string]int64) error {
	muted, err := mutedTags(a.db, u.ID)
	if err != nil {
		return err
	}
	var claimed, send []Post
	for _, p := range posts {
		if firstSeen[p.Slug] < u.CreatedAt {
			continue
		}
		res, err := a.db.Exec(`INSERT OR IGNORE INTO deliveries (user_id, slug) VALUES (?, ?)`, u.ID, p.Slug)
		if err != nil {
			return err
		}
		if n, _ := res.RowsAffected(); n == 1 {
			claimed = append(claimed, p)
			if u.EmailOn && wanted(p, muted) {
				send = append(send, p)
			}
		}
	}
	if len(send) == 0 {
		return nil
	}
	if err := a.mail.Send(a.digest(u, send)); err != nil {
		for _, p := range claimed {
			// The next run sends the post again only if this row is gone.
			if _, derr := a.db.Exec(`DELETE FROM deliveries WHERE user_id = ? AND slug = ?`, u.ID, p.Slug); derr != nil {
				err = errors.Join(err, derr)
			}
		}
		return err
	}
	return nil
}

func (a *app) unsubscribeURL(u user) string {
	return a.cfg.BaseURL + "/account/unsubscribe?t=" + url.QueryEscape(u.UnsubToken)
}

func (a *app) digest(u user, posts []Post) message {
	subject := fmt.Sprintf("%d new posts on Recently Written", len(posts))
	if len(posts) == 1 {
		subject = "New post: " + posts[0].Title
	}
	var b strings.Builder
	b.WriteString("New on Recently Written:\n\n")
	for _, p := range posts {
		fmt.Fprintf(&b, "%s\n%s/%s.html\n\n", p.Title, a.cfg.BaseURL, p.Slug)
	}
	fmt.Fprintf(&b, "-- \nChange your settings: %s/account/\nStop these emails: %s\n",
		a.cfg.BaseURL, a.unsubscribeURL(u))
	return message{To: u.Email, Subject: subject, Body: b.String(), Unsubscribe: a.unsubscribeURL(u)}
}

// notifyLoop runs notify now and then one time in each interval, until
// stop closes.
func (a *app) notifyLoop(interval time.Duration, stop <-chan struct{}) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		if err := a.notify(); err != nil {
			a.log.Printf("notify: %v", err)
		}
		select {
		case <-stop:
			return
		case <-t.C:
		}
	}
}
