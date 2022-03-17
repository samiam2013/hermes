package libmailgun

import (
	"fmt"
	"os"

	mailgun "github.com/mailgun/mailgun-go"
)

// MGEmail implements Sendable for Email
type MGEmail struct {
	FromAddr string
	Subject  string
	TextBody string
	ToAddr   string
	//FromName
	//HTMLBody
	//ReplyTo(name)?
}

func NewMGEmail(to, from, subject, text string) MGEmail {
	return MGEmail{
		FromAddr: from,
		Subject:  subject,
		TextBody: text,
		ToAddr:   to,
	}
}

// Send implements sendTransactional emails with MailGun via the Email type
func (email *MGEmail) Send() error {
	baseURL := os.Getenv("MAILGUN_BASE_URL")
	apiKey := os.Getenv("MAILGUN_API_KEY")
	mg := mailgun.NewMailgun(baseURL, apiKey)
	m := mg.NewMessage(
		email.FromAddr, //? concat e.FromName?
		email.Subject,
		email.TextBody, //? what about html? NewMIMEMessage()?
		email.ToAddr,
	)
	msg, _, err := mg.Send(m)
	if err != nil {
		return fmt.Errorf("failed sending email via sendgrid: %s (%s)",
			msg, err.Error())
	}
	return nil
}
