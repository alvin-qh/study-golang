package smtp

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

const (
	template = "To: %s\r\n" +
		"Subject: %s\r\n" +
		"From: %s <%s>\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		"%s"
)

type SMTP struct {
	Host     string
	Port     int
	Auth     smtp.Auth
	Sender   string
	Nickname string
}

func NewSMTP() (*SMTP, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return nil, err
	}

	return &SMTP{
		Host:     os.Getenv("SERVER"),
		Port:     port,
		Auth:     smtp.PlainAuth("", os.Getenv("ACCOUNT"), os.Getenv("PASSWORD"), os.Getenv("SERVER")),
		Sender:   os.Getenv("SENDER"),
		Nickname: os.Getenv("NICKNAME"),
	}, nil
}

func (s *SMTP) Send(to string, subject string, msg string) error {
	msg = fmt.Sprintf(template, to, subject, s.Nickname, s.Sender, msg)
	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.Host, s.Port),
		s.Auth,
		s.Sender,
		[]string{to},
		[]byte(msg),
	)
}
