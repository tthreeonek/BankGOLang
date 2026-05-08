package service

import (
	"fmt"

	"github.com/go-mail/mail/v2"
	"github.com/sirupsen/logrus"
)

type EmailService struct {
	dialer *mail.Dialer
	from   string
}

func NewEmailService(host string, port int, user, pass string) *EmailService {
	d := mail.NewDialer(host, port, user, pass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	return &EmailService{dialer: d, from: user}
}

func (s *EmailService) SendPaymentNotification(to string, amount float64) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Платеж выполнен")
	m.SetBody("text/html", fmt.Sprintf("<h2>Списание %f RUB</h2>", amount))
	if err := s.dialer.DialAndSend(m); err != nil {
		logrus.Errorf("Email send error: %v", err)
		return err
	}
	logrus.Infof("Sent payment email to %s", to)
	return nil
}
