package libsendinblue

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sendinblue/APIv3-go-library/lib"
	sendinblue "github.com/sendinblue/APIv3-go-library/lib"
	log "github.com/sirupsen/logrus"
)

// BlueEmail struct to be created and sent through library functions
type BlueEmail struct {
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
func NewTextEmail(to, fromName, fromAddr, subject, replyToName, replyToAddr, html string, text []byte) BlueEmail {
	e := BlueEmail{
		To:          to,
		FromName:    fromName,
		FromAddr:    fromAddr,
		Subject:     subject,
		Text:        text,
		HTML:        html,
		ReplyToName: replyToName,
		ReplyToAddr: replyToAddr,
	}
	return e
}

// Send take a SendInBlue API key and creates a new client and sends an Email struct
func (e *BlueEmail) Send() error {
	var ctx context.Context
	apiKey := os.Getenv("SENDINBLUE_API_KEY")
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
		HtmlContent: e.HTML,
		Subject:     e.Subject,
		ReplyTo: &lib.SendSmtpEmailReplyTo{
			Email: e.ReplyToAddr,
			Name:  e.ReplyToName,
		},
	}
	create, httpResp, err := sibClient.TransactionalEmailsApi.SendTransacEmail(ctx, body)
	if err != nil {
		var body []byte
		if body, err = ioutil.ReadAll(httpResp.Body); err != nil {
			log.Error("failed to read response from SIB:", err.Error())
		}
		return fmt.Errorf("error sending email: %s httpResp from SendBlue(): %+v", err.Error(), string(body))
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
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			log.Errorln("Failed reading body of response for error", err.Error())
		}
		log.Errorln("GetAccount Object:", result, " GetAccount Response: ", string(body))
	}
	return
}
