package utils

import (
	"log"
	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "bajuanakanak.reva@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verifikasi Email Anda")
	m.SetBody("text/html", "Klik link berikut untuk verifikasi email Anda: <a href='http://localhost:8000/verify-email?token="+token+"'>Verifikasi Email</a>")

	d := gomail.NewDialer("smtp.gmail.com", 587, "bajuanakanak.reva@gmail.com", "duotiwswtmuidust")

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Gagal mengirim email: %v", err)
		return err
	}
	return nil
}
