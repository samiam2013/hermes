package sendinblue

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sendinblue/APIv3-go-library/lib"
	sendinblue "github.com/sendinblue/APIv3-go-library/lib"
	log "github.com/sirupsen/logrus"
)

// Email struct to be created and sent through library functions
type Email struct {
	To          string
	FromName    string
	FromAddr    string
	Subject     string
	ReplyToName string
	ReplyToAddr string
	Text        []byte
	HTML        string // currently unused
}

// NewTextEmail returns a simple text Email struct with attached Send.*() funcs
func NewTextEmail(to, fromName, fromAddr, subject, replyToName, replyToAddr string, text []byte) Email {
	e := Email{
		To:          to,
		FromName:    fromName,
		FromAddr:    fromAddr,
		Subject:     subject,
		Text:        text,
		HTML:        "",
		ReplyToName: replyToName,
		ReplyToAddr: replyToAddr,
	}
	return e
}

// SendBlue take a SendInBlue API key and creates a new client and sends an Email struct
func (e *Email) SendBlue(apiKey string) error {
	var ctx context.Context
	sibClient := newBlueClient(ctx, apiKey)
	textcontent := string(e.Text) // TODO what if this text is really REALLY big?
	body := lib.SendSmtpEmail{
		Sender: &lib.SendSmtpEmailSender{
			Name:  e.FromName,
			Email: e.FromAddr,
		},
		To: []lib.SendSmtpEmailTo{
			{Email: e.To}},
		TextContent: textcontent, // using text here, not html to help prevent XSS via email
		Subject:     e.Subject,
		ReplyTo: &lib.SendSmtpEmailReplyTo{
			Email: e.ReplyToAddr,
			Name:  e.ReplyToName,
		},
	}
	create, httpResp, err := sibClient.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		return fmt.Errorf("error sending email: %s httpResp from SendBlue(): %+v", err.Error(), httpResp)
	}
	if httpResp.StatusCode != http.StatusCreated {
		return fmt.Errorf("created msg with id %v, *non*-accepted http code: %d",
			create.MessageId, httpResp.StatusCode)
	}
	return nil
}

// newSIBClient opens and returns a new sendinblue APIClient
func newBlueClient(ctx context.Context, apiKey string) (sib *sendinblue.APIClient) {

	sibCfg := sendinblue.NewConfiguration()
	//Configure API key authorization: api-key
	sibCfg.AddDefaultHeader("api-key", apiKey)
	//Configure API key authorization: partner-key
	sibCfg.AddDefaultHeader("partner-key", apiKey)

	sib = sendinblue.NewAPIClient(sibCfg)
	result, resp, err := sib.AccountApi.GetAccount(ctx)
	if err != nil {
		log.Errorf("Error when calling AccountApi->get_account: %s ", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorln("GetAccount Object:", result, " GetAccount Response: ", resp)
	}
	return
}
