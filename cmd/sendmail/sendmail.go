package main

import "fmt"

func main() {
	fmt.Println("sendmail binary not implemented")

	// parse the environment variables
	//	fall back onto the .env file in this folder (cmd/sendmail)

	// take in the first argument and if it doesn't parse
	//	as a valid email address panic on the user

	// take in the first line of input and if it doesn't
	//	start with '/subject\s:/' panic on the user

	// take in as many as some hundred lines into a buffer
	//	utnil a single '.' [period] on a line or ctrl + d

	// try each platform with credentials parsed from .env
	//	or the environment until one works

	// print out the api response from said platform

}
