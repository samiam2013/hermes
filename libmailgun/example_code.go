package libmailgun

import mailgun "github.com/mailgun/mailgun-go"

func SendSimpleMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"hermes <no-reply@goon.villas>",
		"Hello",
		"Testing some Mailgun awesomeness!",
		"sam@myres.dev",
	)
	_, id, err := mg.Send(m)
	return id, err
}
