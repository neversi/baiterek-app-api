package service

import (
	"errors"
	"io"
	"log"
	"time"

	"github.com/emersion/go-smtp"
)

// The Backend implements SMTP server methods.
type Backend struct{}

const (
	susername = "mlewis"
	spassword = "mlewis"
)

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username == susername && password == spassword {
		return &Session{}, nil
	}

	return nil, smtp.ErrAuthRequired
}

func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct{}

func (s *Session) AuthPlain(username, password string) error {
	if username != "username" || password != "password" {
		return errors.New("invalid username or password")
	}
	return nil
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	log.Println("Data:", string(b))

	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func StartSMTPServer(port string) error {
	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = port
	s.Domain = "0.0.0.0"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
