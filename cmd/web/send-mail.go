// listen to the channel

package main

import (
	"log"
	"time"

	"github.com/valorvig/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

/*
// this function needs to be called every time I want to send a message
func ListenForMail() {
	m := <-app.MailChan
}
*/

func ListenForMail() { // this func needs to be called somewherer, e.g. in main func
	// start off in the background using go routine
	go func() {
		// listen all the time for incoming data
		for {
			msg := <-app.MailChan
			sendMsg(msg) // try testing by running the app after putting this line
		}
	}()
}

// return nothing, only send email
func sendMsg(m models.MailData) {
	// Define a server
	server := mail.NewSMTPClient()
	server.Host = "localhost" // Where is the server? - we're using our dummy mail server
	server.Port = 1025        // In most cases, real mail servers use Port 587 or 465
	server.KeepAlive = false  // don't want the mail active all the time
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// these three are needed in production ------
	// server.Username
	// server.Password
	// server.Encryption

	// Define a client
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	// set from, to, and subject
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	// Separate the body - we don't want to manually construct an email
	// Requires a type of content and the content
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}
}
