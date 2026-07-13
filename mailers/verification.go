package mailers

import (
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
)

func SendVerification(email string, link string) error {
	m := mail.NewMessage()

	m.Subject = "Verify Your Account"
	m.From = "no-reply@example.com"
	m.To = []string{email}

	err := m.AddBody(r.HTML("mail/verification.plush.html"), render.Data{
		"link": link,
	})
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
