package libsendgrid

import (
	"fmt"
	"net/http"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

// GridEmail contains basic email data to, from, etc. and can be easily extended
type GridEmail struct {
	FromAddr string
	FromName string
	ToAddr   string
	Subject  string
	TextBody string
	HTML     string
}

// Send reciever function calls sendgrid API for a transactional email
func (email *GridEmail) Send(sendGridKey string) error {
	var from *mail.Email
	if email.FromAddr == "" {
		return fmt.Errorf("cannot send without email.FromAddr defined")
	}
	from = mail.NewEmail(email.FromName, email.FromAddr)
	sendTo := mail.NewEmail("", email.ToAddr)

	var html string
	if email.HTML != "" {
		// TODO sanitize with bluemonday
		html = email.HTML
	} else {
		html = "<html><pre>" + email.TextBody + "</pre></html>"
	}

	message := mail.NewSingleEmail(from, email.Subject, sendTo, email.TextBody, html)
	client := sendgrid.NewSendClient(sendGridKey)
	response, err := client.Send(message)
	if err != nil {
		log.Warn("sendgrid email failed: ", response.StatusCode,
			response.Body, response.Headers)
		return err
	} else {
		if response.StatusCode == http.StatusAccepted {
			log.Info("email was accepted by sendgrid")
		} else {
			log.Warn("email was accepted but with status code:",
				response.StatusCode, "body:", response.Body)
		}
	}
	return nil
}
