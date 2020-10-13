package logic

import (
	"log"

	mail "gopkg.in/mail.v2"
)

type MailTask struct {
	from    string
	to      string
	content string
}

func (mt *MailTask) Send(mailDialer *mail.Dialer) {
	//TODO Implement error hadleling
	message := mail.NewMessage()
	//Should be on env viriables
	message.SetHeader("From", mt.from)

	message.SetHeader("To", mt.to)
	message.SetBody("text/plain", mt.content)

	//Sending Email
	err := mailDialer.DialAndSend(message)

	if err != nil {
		log.Println(err)
	}
}
