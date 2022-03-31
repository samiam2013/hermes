package hermes

import (
	"fmt"
	"log"
	"os"

	"github.com/microcosm-cc/bluemonday"
	"github.com/samiam2013/hermes/libmailgun"
	"github.com/samiam2013/hermes/libsendgrid"
	"github.com/samiam2013/hermes/libsendinblue"
)

// Email defineds a data structure for a single message from one to one person independent of platform
type Email struct {
	ToName      string
	ToAddr      string
	TextBody    string
	Subject     string
	ReplyToName string
	ReplyToAddr string
	HTMLBody    string
	FromName    string
	FromAddr    string
	credentials map[uint]map[string]string
}

// Sendable interface enforces structure across sub-libraries
type Sendable interface {
	Send() error
}

var _ Sendable = &libmailgun.MGEmail{}
var _ Sendable = &libsendinblue.BlueEmail{}
var _ Sendable = &libsendgrid.GridEmail{}

// constant uint values for the choices of platform
const (
	SendInBlue = iota
	SendGrid
	MailGun
)

// EnvVarNames stores the set of required .env keys required
//	by each platform when sending
var EnvVarNames = map[uint]map[string]string{
	SendGrid: {
		"key":    "SENDGRID_API_KEY",
		"sender": "SENDGRID_SENDER",
	},
	SendInBlue: {
		"key":    "SENDINBLUE_API_KEY",
		"sender": "SENDINBLUE_SENDER",
	},
	MailGun: {
		"key":     "MAILGUN_API_KEY",
		"baseurl": "MAILGUN_BASE_URL",
	},
}

// NewTransactional looks at the environment and returns a sendable email or an error
func NewTransactional() (Email, error) {
	// try parsing env for a set of platform variables
	creds, platform, err := parseCreds()
	if err != nil {
		return Email{}, fmt.Errorf("failed to parseCreds() for NewTransactional(): %s", err.Error())
	}
	return NewTransactionalWithEnv(creds, platform), nil
}

// NewTransactionalWithEnv can be called with the environment variables rather than parsing them with os.GetEnv
func NewTransactionalWithEnv(creds map[string]string, platform uint) Email {
	// hand back an empty email struct with platform creds list
	new := Email{
		ToName:      "",
		ToAddr:      "",
		TextBody:    "",
		Subject:     "",
		ReplyToName: "",
		ReplyToAddr: "",
		HTMLBody:    "",
		FromName:    "",
		FromAddr:    creds[EnvVarNames[platform]["sender"]],
		credentials: map[uint]map[string]string{platform: creds},
	}
	return new
}

func parseCreds() (map[string]string, uint, error) {
	// parse .env from this folder?
	credentials := map[string]string{}
	foundSet := false
	platform := uint(len(EnvVarNames)) // will cause an error if indexed
	for platformID, requiredSet := range EnvVarNames {
		satisfied := true
		for _, envVar := range requiredSet {
			if os.Getenv(envVar) == "" {
				fmt.Println("checking for envVar:", envVar)
				satisfied = false
			}
		}
		if satisfied {
			foundSet = true
			for _, found := range requiredSet {
				credentials[found] = os.Getenv(found)
			}
			// if a platform wasn't already found (platform has init value) set it
			if platform == uint(len(EnvVarNames)) {
				platform = platformID
			} else {
				// what to do if a platform was already found ?
				continue
			}
		}
	}
	if !foundSet {
		return nil, uint(len(EnvVarNames) + 1), // platform id will cause a panic if accessed blindly (intended effect)
			fmt.Errorf("did not find a set of credentials in the environment")
	}
	return credentials, platform, nil

}

// Send a platform ambiguous email structure :D
func (e *Email) Send() (err error) {
	if len(e.credentials) == 0 {
		return fmt.Errorf("no credentials set on email! (do you need hermes.NewTransactional()?)")
	}
	for platformID := range e.credentials {
		switch platformID {
		case SendGrid:
			if err = e.sendSendGrid(); err != nil {
				log.Printf("sendgrid send attempt failed: %s\n", err.Error())
				continue
			}
			return
		case SendInBlue:
			if err = e.sendSendInBlue(); err != nil {
				log.Printf("sendblue send attempt failed: %s\n", err.Error())
				continue
			}
			return
		case MailGun:
			if err = e.sendMailGun(); err != nil {
				log.Printf("sendgrid send attempt failed: %s\n", err.Error())
				continue
			}
			return

		default:
			return fmt.Errorf("platform %d not resolved", platformID)
		}
	}
	return nil
}

func (e *Email) sendSendGrid() error {
	apiSenderIdx := EnvVarNames[SendGrid]["sender"]
	sgEmail, err := libsendgrid.NewGridEmail(
		e.credentials[SendGrid][apiSenderIdx],
		e.ToAddr,
		e.Subject,
		e.TextBody,
		sanitizeHTML(e.HTMLBody))
	//e.FromAddr
	//e.FromName,
	if err != nil {
		log.Fatalf("failed to create sendgrid message: %s", err.Error())
	}
	//ToName: e.ToName,
	//ReplyToName: e.ReplyToName
	//ReplyToAddr: e.ReplyToAddr
	return sgEmail.Send()
}

func (e *Email) sendSendInBlue() error {
	blueEmail := libsendinblue.NewTextEmail(
		e.ToAddr,
		e.FromName,
		e.FromAddr,
		e.Subject,
		e.ReplyToName,
		e.ReplyToAddr,
		sanitizeHTML(e.HTMLBody),
		[]byte(e.TextBody),
		// e.ToName, e.HTML still available
	)
	return blueEmail.Send()

}

func (e *Email) sendMailGun() error {
	newEmail := libmailgun.NewMGEmail(
		e.ToAddr,
		e.FromAddr,
		e.Subject,
		e.TextBody,
		sanitizeHTML(e.HTMLBody),
	)
	return newEmail.Send()
}

func sanitizeHTML(input string) string {
	p := bluemonday.UGCPolicy()
	output := p.Sanitize(input)
	return output
}
