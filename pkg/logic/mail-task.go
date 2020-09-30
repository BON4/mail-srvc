package logic

import (
	mail "gopkg.in/mail.v2"
)

type MailTask struct {
	from    string
	to      string
	content string
}

func (mt *MailTask) Send(mailDialer *mail.Dialer) error {
	message := mail.NewMessage()
	//Should be on env viriables
	message.SetHeader("From", "from@gmail.com")

	message.SetHeader("To", mt.to)
	message.SetBody("text/plain", mt.content)

	//Sending Email
	if err := mailDialer.DialAndSend(message); err != nil {
		return err
	}

	return nil
}
