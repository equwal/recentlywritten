package main

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	_ "modernc.org/sqlite"
)

// The users table holds only addresses that a sign-in link verified.
// login_tokens and sessions hold SHA-256 hashes of their tokens, so a copy
// of the database does not let a person sign in. feed_token and
// unsub_token are stored as plain text: the site must show them again,
// and they give access only to public posts and to the unsubscribe step.
const schema = `
CREATE TABLE IF NOT EXISTS users (
	id          INTEGER PRIMARY KEY,
	email       TEXT NOT NULL UNIQUE,
	created_at  INTEGER NOT NULL,
	email_on    INTEGER NOT NULL DEFAULT 1,
	feed_token  TEXT NOT NULL UNIQUE,
	unsub_token TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS muted_tags (
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	tag     TEXT NOT NULL,
	PRIMARY KEY (user_id, tag)
);
CREATE TABLE IF NOT EXISTS login_tokens (
	hash       TEXT PRIMARY KEY,
	email      TEXT NOT NULL,
	created_at INTEGER NOT NULL,
	expires_at INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS sessions (
	hash       TEXT PRIMARY KEY,
	user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	expires_at INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS posts (
	slug       TEXT PRIMARY KEY,
	first_seen INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS deliveries (
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	slug    TEXT NOT NULL,
	PRIMARY KEY (user_id, slug)
);
`

const (
	loginTTL     = 15 * time.Minute
	loginSpacing = 2 * time.Minute // minimum time between two links to one address
	sessionTTL   = 30 * 24 * time.Hour
)

var errNotFound = errors.New("not found")

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	// One connection serializes all writes. The site is small, and this
	// removes the need to handle SQLITE_BUSY. It also keeps a :memory:
	// database the same for all queries in tests.
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// newToken returns a random token with 128 bits of entropy.
func newToken() string { return rand.Text() }

func hashToken(t string) string {
	s := sha256.Sum256([]byte(t))
	return hex.EncodeToString(s[:])
}

type user struct {
	ID         int64
	Email      string
	CreatedAt  int64
	EmailOn    bool
	FeedToken  string
	UnsubToken string
}

const userColumns = `id, email, created_at, email_on, feed_token, unsub_token`

func scanUser(row *sql.Row) (user, error) {
	var u user
	err := row.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.EmailOn, &u.FeedToken, &u.UnsubToken)
	if errors.Is(err, sql.ErrNoRows) {
		err = errNotFound
	}
	return u, err
}

// createLoginToken makes a sign-in token for email. It returns an empty
// token and no error if a token for email was made less than loginSpacing
// ago. This stops a person who uses the form to send many emails to one
// address.
func createLoginToken(db *sql.DB, email string, now time.Time) (string, error) {
	if _, err := db.Exec(`DELETE FROM login_tokens WHERE expires_at <= ?`, now.Unix()); err != nil {
		return "", err
	}
	var recent int
	err := db.QueryRow(`SELECT count(*) FROM login_tokens WHERE email = ? AND created_at > ?`,
		email, now.Add(-loginSpacing).Unix()).Scan(&recent)
	if err != nil || recent > 0 {
		return "", err
	}
	t := newToken()
	_, err = db.Exec(`INSERT INTO login_tokens (hash, email, created_at, expires_at) VALUES (?, ?, ?, ?)`,
		hashToken(t), email, now.Unix(), now.Add(loginTTL).Unix())
	if err != nil {
		return "", err
	}
	return t, nil
}

// consumeLoginToken removes the token and returns the user for its
// address. It makes the user if the address has no account. A token works
// one time only.
func consumeLoginToken(db *sql.DB, token string, now time.Time) (user, error) {
	var email string
	err := db.QueryRow(`DELETE FROM login_tokens WHERE hash = ? AND expires_at > ? RETURNING email`,
		hashToken(token), now.Unix()).Scan(&email)
	if errors.Is(err, sql.ErrNoRows) {
		return user{}, errNotFound
	}
	if err != nil {
		return user{}, err
	}
	_, err = db.Exec(`INSERT INTO users (email, created_at, feed_token, unsub_token) VALUES (?, ?, ?, ?)
		ON CONFLICT (email) DO NOTHING`, email, now.Unix(), newToken(), newToken())
	if err != nil {
		return user{}, err
	}
	return scanUser(db.QueryRow(`SELECT `+userColumns+` FROM users WHERE email = ?`, email))
}

func createSession(db *sql.DB, userID int64, now time.Time) (string, error) {
	if _, err := db.Exec(`DELETE FROM sessions WHERE expires_at <= ?`, now.Unix()); err != nil {
		return "", err
	}
	t := newToken()
	_, err := db.Exec(`INSERT INTO sessions (hash, user_id, expires_at) VALUES (?, ?, ?)`,
		hashToken(t), userID, now.Add(sessionTTL).Unix())
	return t, err
}

func sessionUser(db *sql.DB, token string, now time.Time) (user, error) {
	return scanUser(db.QueryRow(`SELECT u.id, u.email, u.created_at, u.email_on, u.feed_token, u.unsub_token
		FROM sessions s JOIN users u ON u.id = s.user_id
		WHERE s.hash = ? AND s.expires_at > ?`, hashToken(token), now.Unix()))
}

func deleteSession(db *sql.DB, token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE hash = ?`, hashToken(token))
	return err
}

func userByFeedToken(db *sql.DB, token string) (user, error) {
	return scanUser(db.QueryRow(`SELECT `+userColumns+` FROM users WHERE feed_token = ?`, token))
}

func mutedTags(db *sql.DB, userID int64) (map[string]bool, error) {
	rows, err := db.Query(`SELECT tag FROM muted_tags WHERE user_id = ?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	muted := map[string]bool{}
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		muted[t] = true
	}
	return muted, rows.Err()
}

// saveSettings sets the email switch and replaces the muted tags.
func saveSettings(db *sql.DB, userID int64, emailOn bool, muted []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE users SET email_on = ? WHERE id = ?`, emailOn, userID); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM muted_tags WHERE user_id = ?`, userID); err != nil {
		return err
	}
	for _, t := range muted {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO muted_tags (user_id, tag) VALUES (?, ?)`, userID, t); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// rotateFeedToken gives the user a new feed token. The old feed URL then
// stops working.
func rotateFeedToken(db *sql.DB, userID int64) error {
	_, err := db.Exec(`UPDATE users SET feed_token = ? WHERE id = ?`, newToken(), userID)
	return err
}

// unsubscribe turns off email for the user with the unsubscribe token. It
// returns errNotFound for an unknown token.
func unsubscribe(db *sql.DB, token string) error {
	res, err := db.Exec(`UPDATE users SET email_on = 0 WHERE unsub_token = ?`, token)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errNotFound
	}
	return nil
}

// deleteUser removes the user. The foreign keys remove the sessions,
// muted tags and delivery records of the user.
func deleteUser(db *sql.DB, userID int64) error {
	_, err := db.Exec(`DELETE FROM users WHERE id = ?`, userID)
	return err
}
