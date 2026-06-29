package mailers

import (
	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
)

func SendAdminPasswordReset(email string, link string) error {
	m := mail.NewMessage()

	m.Subject = "Admin Password Reset"
	m.From = "no-reply@example.com"
	m.To = []string{email}

	err := m.AddBody(r.HTML("mail/admin_password_reset.plush.html"), render.Data{
		"link": link,
	})
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
