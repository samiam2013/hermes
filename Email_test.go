package hermes

import (
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func TestNewTransactionWithEnv(t *testing.T) {
	type args struct {
		creds    map[string]string
		platform uint
	}

	sibCreds := map[string]string{
		"SENDINBLUE_API_KEY": (func() string {
			err := godotenv.Load("cmd/sendmail/.env")
			if err != nil {
				panic("test needs cmd/sendmail/.env")
			}
			return os.Getenv("SENDINBLUE_API_KEY")
		}()),
		"SENDINBLUE_SENDER": (func() string {
			err := godotenv.Load("cmd/sendmail/.env")
			if err != nil {
				panic("test needs cmd/sendmail/.env")
			}
			return os.Getenv("SENDINBLUE_SENDER")
		}()),
	}

	tests := []struct {
		name string
		args args
		want Email
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				creds:    sibCreds,
				platform: 0,
			},
			want: Email{
				FromAddr:    sibCreds["SENDINBLUE_SENDER"],
				FromName:    "",
				ToAddr:      "",
				ToName:      "",
				Subject:     "",
				ReplyToName: "",
				ReplyToAddr: "",
				TextBody:    "",
				HTMLBody:    "",
				credentials: credentials{
					set:      true,
					platform: 0,
					list:     sibCreds,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//log.Printf("creds: %+v", tt.args.creds)
			if got := NewTransactionalWithEnv(tt.args.creds, tt.args.platform); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTransactionWithEnv() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
