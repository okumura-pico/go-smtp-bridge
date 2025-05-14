package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

var auth = sasl.NewPlainClient("", getEnv("SMTP_BRIDGE_USER"), getEnv("SMTP_BRIDGE_PASSWORD"))
var addr = getEnv("SMTP_BRIDGE_HOST") + ":" + getEnv("SMTP_BRIDGE_PORT")

// The Backend implements SMTP server methods.
type Backend struct{}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

type Session struct {
	From   string
	RcptTo string
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	log.Println("Rcpt to:", to)
	s.RcptTo = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	return s.ForwardMail(r)
}

func (s *Session) Reset() {
}

func (s *Session) Logout() error {
	return nil
}

func (s *Session) ForwardMail(r io.Reader) error {
	to := []string{s.RcptTo}

	return smtp.SendMail(addr, auth, s.From, to, r)
}

func getEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("Define environment variable: %s", key)
	}

	return value
}

func main() {
	be := &Backend{}

	s := smtp.NewServer(be)

	if len(os.Args) == 1 {
		s.Addr = "localhost:1025"
	} else {
		s.Addr = os.Args[1]
	}

	s.Domain = "localhost"
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 30 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
