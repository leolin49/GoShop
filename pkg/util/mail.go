package util

import "net/smtp"

type EmailInfo struct {
	Sender  string
	Recver  string
	Subject string
	Body    string
}

// TODO
func SendEmail(mail EmailInfo) {
	smtpHost := "smtp.example.com"
	smtpPort := "587"
	pwd := "123456"

	to := []string{mail.Recver}
	auth := smtp.PlainAuth("", mail.Sender, pwd, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, mail.Sender, to, []byte(mail.Subject+"\n"+mail.Body))
	if err != nil {
		return
	}
	// success
}
