package hermes

import (
	"fmt"
	"os"

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
	new := Email{}
	new.FromAddr = creds[EnvVarNames[platform]["sender"]]
	new.credentials = make(map[uint]map[string]string)
	new.credentials[platform] = creds
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
func (e *Email) Send() error {
	if len(e.credentials) == 0 {
		return fmt.Errorf("no credentials set on email! (do you need hermes.NewTransactional()?)")
	}
	for platformID := range e.credentials {
		switch platformID {
		case SendGrid:
			return e.sendSendGrid()
		case SendInBlue:
			return e.sendSendInBlue()
		case MailGun:
			return e.sendMailGun()

		default:
			return fmt.Errorf("platform not resolved")
		}
	}
	return nil
}

func (e *Email) sendSendGrid() error {
	apiSenderIdx := EnvVarNames[SendGrid]["sender"]
	sgEmail := libsendgrid.GridEmail{
		FromAddr: e.credentials[SendGrid][apiSenderIdx], //or e.FromAddr?
		FromName: e.FromName,
		ToAddr:   e.ToAddr,
		//ToName: e.ToName,
		//ReplyToName: e.ReplyToName
		//ReplyToAddr: e.ReplyToAddr
		Subject:  e.Subject,
		TextBody: e.TextBody,
		HTML:     e.HTMLBody,
	}
	return sgEmail.Send()
}

func (e *Email) sendSendInBlue() error {
	newEmail := libsendinblue.NewTextEmail(
		e.ToAddr,
		e.FromName,
		e.FromAddr,
		e.Subject,
		e.ReplyToName,
		e.ReplyToAddr,
		[]byte(e.TextBody))
	// e.ToName, e.HTML still available
	return newEmail.Send()

}

func (e *Email) sendMailGun() error {
	newEmail := libmailgun.NewMGEmail(
		e.ToAddr,
		e.FromAddr,
		e.Subject,
		e.TextBody,
	)
	return newEmail.Send()
}
