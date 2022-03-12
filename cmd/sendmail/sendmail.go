package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/samiam2013/hermes"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("failed to read environment file at ./.env: %s", err.Error())
	}

	email, err := hermes.NewTransactional()
	if err != nil {
		panic(err.Error())
	}
	// parse the environment variables if they appear as a flag
	//	fall back onto the .env file in this folder (cmd/sendmail)
	email.ToAddr = "sam@myres.dev"
	email.FromName = "hermes roboto"
	email.FromAddr = "sam@myres.dev"
	// take in the first argument and if it doesn't parse
	//	as a valid email address panic on the user

	// take in the first line of input and if it doesn't
	//	start with '/subject\s:/' panic on the user
	email.Subject = "test from hermes"

	// take in as many as some hundred lines into a buffer
	//	utnil a single '.' [period] on a line or ctrl + d
	email.TextBody = "this is the 'text' body of the email."

	// try each platform with credentials parsed from .env
	//	or the environment until one works

	email.FromName = "hermes"

	err = email.Send()
	if err != nil {
		panic(err.Error())
	}
}
