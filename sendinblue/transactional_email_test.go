package sendinblue

import (
	"reflect"
	"testing"
)

func TestEmail_SendBlue(t *testing.T) {
	type fields struct {
		To          string
		FromName    string
		FromAddr    string
		Subject     string
		ReplyToName string
		ReplyToAddr string
		Text        []byte
		HTML        string
	}
	type args struct {
		apiKey string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add secret mechanism for testing
		// against a separate API key
		{
			name: "unhappy path",
			fields: fields{
				To:          "OtherGuy@Myres.dev",
				FromName:    "Go Automated Test",
				FromAddr:    "Sam@Myres.dev",
				Subject:     "Go hermes SendBlue() test",
				ReplyToName: "Sam Myres",
				ReplyToAddr: "Sam@Myres.dev",
				Text:        []byte("this is body text for the test"),
				HTML:        "",
			},
			args: args{
				// I'm sorry for this.
				apiKey: func() string {
					// err := godotenv.Load()
					// if err != nil {
					// 	t.Errorf("error loading .env file for test: %s", err.Error())
					// }
					// key, ok := os.LookupEnv("SIB_APIKEY")
					// if !ok {
					// 	t.Errorf("SIB_APIKEY not in env vars for test")
					// }
					// return key
					return "bad_key"
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &BlueEmail{
				To:          tt.fields.To,
				FromName:    tt.fields.FromName,
				FromAddr:    tt.fields.FromAddr,
				Subject:     tt.fields.Subject,
				ReplyToName: tt.fields.ReplyToName,
				ReplyToAddr: tt.fields.ReplyToAddr,
				Text:        tt.fields.Text,
				HTML:        tt.fields.HTML,
			}
			if err := e.Send(tt.args.apiKey); (err != nil) != tt.wantErr {
				t.Errorf("Email.SendBlue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTextEmail(t *testing.T) {
	type args struct {
		to          string
		fromName    string
		fromAddr    string
		subject     string
		replyToName string
		replyToAddr string
		text        []byte
	}
	tests := []struct {
		name string
		args args
		want BlueEmail
	}{
		// TODO: Add test cases.
	{
		name: "",
		args: args{
			to:          "to@place.tld",
			fromName:    "automated sender",
			fromAddr:    "automated@place.tld",
			subject:     "some subject line",
			replyToName: "reply toguy",
			replyToAddr: "replyto@place.tld",
			text:        []byte("some text email"),
		},
		want: BlueEmail{
			To:          "to@place.tld",
			FromName:    "automated sender",
			FromAddr:    "automated@place.tld",
			Subject:     "some subject line",
			ReplyToName: "reply toguy",
			ReplyToAddr: "replyto@place.tld",
			Text:        []byte("some text email"),
			HTML:        "",
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTextEmail(tt.args.to, tt.args.fromName, tt.args.fromAddr, tt.args.subject, tt.args.replyToName, tt.args.replyToAddr, tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTextEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
