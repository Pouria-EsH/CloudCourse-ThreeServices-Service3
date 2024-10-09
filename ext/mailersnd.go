package ext

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailersend/mailersend-go"
)

type MailerSend struct {
	from mailersend.From
	ms   *mailersend.Mailersend
}

func NewMailerSend(apikey string, fromemail string, fromname string) *MailerSend {
	return &MailerSend{
		from: mailersend.From{
			Name:  fromname,
			Email: fromemail,
		},
		ms: mailersend.NewMailersend(apikey),
	}
}

func (m MailerSend) Send(email Email) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	recipients := []mailersend.Recipient{
		{
			Name:  email.RecipientName,
			Email: email.RecipientMail,
		},
	}

	message := m.ms.Email.NewMessage()

	message.SetFrom(m.from)
	message.SetRecipients(recipients)
	message.SetSubject(email.Subject)
	message.SetText(email.Message)

	log.Println("sending email to ", email.RecipientMail)
	_, err := m.ms.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("error in sending email: %v", err)
	}

	log.Println("email sent")
	return nil
}
