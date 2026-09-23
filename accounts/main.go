// Command accounts is the sign-in, email, and private feed service for
// recentlywritten.com. It serves the paths under /account/ behind the web
// server that serves the static site, and it reads the post data that
// build.sh writes to postdata/ in the site directory.
//
// All settings and secrets come from environment variables. See
// .env.example.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/mail"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type config struct {
	BaseURL      string // for example https://recentlywritten.com, with no final slash
	SiteDir      string // the web root that build.sh deploys
	DBPath       string
	Listen       string
	MailFrom     *mail.Address
	SMTPAddr     string
	SMTPUsername string
	SMTPPassword string
	Interval     time.Duration
}

// loadConfig reads the configuration from getenv. In dev mode the mail
// settings are not necessary. The error names each missing variable.
func loadConfig(getenv func(string) string, dev bool) (config, error) {
	c := config{
		BaseURL:      strings.TrimRight(getenv("RW_BASE_URL"), "/"),
		SiteDir:      getenv("RW_SITE_DIR"),
		DBPath:       getenv("RW_DB_PATH"),
		Listen:       getenv("RW_LISTEN"),
		SMTPUsername: getenv("SMTP_USERNAME"),
		SMTPPassword: getenv("SMTP_PASSWORD"),
		Interval:     15 * time.Minute,
	}
	if c.Listen == "" {
		c.Listen = "127.0.0.1:8081"
	}
	required := []string{"RW_BASE_URL", "RW_SITE_DIR", "RW_DB_PATH"}
	if !dev {
		required = append(required, "RW_MAIL_FROM", "SMTP_HOST")
	}
	var missing []string
	for _, k := range required {
		if getenv(k) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return c, fmt.Errorf("set these environment variables: %s", strings.Join(missing, ", "))
	}
	if host := getenv("SMTP_HOST"); host != "" {
		port := getenv("SMTP_PORT")
		if port == "" {
			port = "587"
		}
		c.SMTPAddr = net.JoinHostPort(host, port)
	}
	if from := getenv("RW_MAIL_FROM"); from != "" {
		addr, err := mail.ParseAddress(from)
		if err != nil {
			return c, fmt.Errorf("RW_MAIL_FROM: %w", err)
		}
		c.MailFrom = addr
	}
	if s := getenv("RW_NOTIFY_INTERVAL"); s != "" {
		d, err := time.ParseDuration(s)
		if err != nil || d <= 0 {
			return c, fmt.Errorf("RW_NOTIFY_INTERVAL: not a positive duration such as 15m")
		}
		c.Interval = d
	}
	return c, nil
}

func main() {
	dev := flag.Bool("dev", false, "write email to standard error, not to SMTP (development only)")
	flag.Parse()
	logger := log.New(os.Stderr, "", log.LstdFlags)

	cfg, err := loadConfig(os.Getenv, *dev)
	if err != nil {
		logger.Fatal(err)
	}
	db, err := openDB(cfg.DBPath)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	var m mailer = logMailer{W: os.Stderr}
	if !*dev {
		m = smtpMailer{Addr: cfg.SMTPAddr, Username: cfg.SMTPUsername, Password: cfg.SMTPPassword, From: cfg.MailFrom}
	}
	a := &app{cfg: cfg, db: db, mail: m, now: time.Now, log: logger}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go a.notifyLoop(cfg.Interval, ctx.Done())

	srv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           a.routes(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
	}
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(shutdown)
	}()
	logger.Printf("listening on %s", cfg.Listen)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal(err)
	}
}
