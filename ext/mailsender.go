package ext

type MailSender interface {
	Send(mail Email) error
}

type Email struct {
	RecipientMail string
	RecipientName string
	Subject       string
	Message       string
}
