package hermes

import (
	"fmt"
	"os"
)

// Email defineds a data structure for a single message from one to one person independent of platform
type Email struct {
	FromAddr    string
	FromName    string
	ToAddr      string
	ToName      string
	Subject     string
	ReplyToName string
	ReplyToAddr string
	TextBody    string
	HTMLBody    string
	credentials credentials
}

// credentials provides a way to store the credentials to send an email alongside the email's data
type credentials struct {
	set      bool
	platform uint
	list     map[string]string
}

// choices of platform
const (
	SendInBlue = iota
	SendGrid
)

var requiredVars = map[uint]map[string]string{
	SendGrid: {
		"key":    "SENDGRID_API_KEY",
		"sender": "SENDGRID_SENDER",
	},
	SendInBlue: {
		"key":    "SENDINBLUE_API_KEY",
		"sender": "SENDINBLUE_SENDER",
	},
}

// NewTransactional looks at the environment and returns a sendable email or an error
func NewTransactional() (Email, error) {
	fmt.Printf("NewTrasactional() not yet implemented")
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
	new.FromAddr = creds[requiredVars[platform]["sender"]]
	new.credentials = credentials{
		set:      true,
		platform: platform,
		list:     creds,
	}
	return new
}

func parseCreds() (map[string]string, uint, error) {
	// parse .env from this folder?
	credentials := map[string]string{}
	foundSet := false
	platform := uint(0)
	for platformID, requiredSet := range requiredVars {
		satisfied := true
		for _, envVar := range requiredSet {
			if os.Getenv(envVar) == "" {
				satisfied = false
			}
		}
		if satisfied {
			foundSet = true
			platform = platformID
			for _, found := range requiredSet {
				credentials[found] = os.Getenv(found)
			}
			break
		}

	}
	if !foundSet {
		return nil, uint(len(requiredVars) + 1), // platform id will cause a panic if accessed blindly (intended effect)
			fmt.Errorf("did not find a set of credentials in the environment")
	}
	return credentials, platform, nil

}

// Send a platform ambiguous email structure :D
func (e *Email) Send() error {
	return fmt.Errorf("not implemented")
}
