package gmail

import (
	"binadesa2020-backend/lib/variable"
	"log"
	"net/smtp"
)

// Email structure
type Email struct {
	From    string
	To      string
	Subject string
	Body    string
}

// Send Email
func (e *Email) Send() error {
	config := variable.GmailConfig
	e.From = config.Email

	msg := "From: " + e.From + "\n" +
		"To: " + e.To + "\n" +
		"Subject: " + e.Subject + "\n\n" +
		e.Body

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		smtp.PlainAuth("", e.From, config.Password, "smtp.gmail.com"),
		e.From,
		[]string{e.To},
		[]byte(msg),
	)
	if err != nil {
		log.Printf("Smtp error to %s : %v", e.To, err)
		return err
	}
	return nil
}

// SendReceiveSubmission via email
func (e *Email) SendReceiveSubmission(typeSubmission string) error {
	e.Subject = "Pengajuan Diterima"
	e.Body = "pengajuan " + typeSubmission + " anda telah diterima oleh sistem kami."

	return e.Send()
}

// SendCompleteSubmission via email
func (e *Email) SendCompleteSubmission(typeSubmission string) error {
	e.Subject = "Pengajuan Telah Selesai"
	e.Body = "pengajuan " + typeSubmission + " anda telah selesai diproses."

	return e.Send()
}
