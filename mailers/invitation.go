package mailers

import (
	"fmt"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
)

func SendInvitation(email string, organization string) error {
	m := mail.NewMessage()

	m.Subject = fmt.Sprintf("You have been added to %s", organization)
	m.From = "no-reply@example.com"
	m.To = []string{email}

	err := m.AddBody(r.HTML("mail/invitation.plush.html"), render.Data{
		"organization": organization,
	})
	if err != nil {
		return err
	}
	return smtp.Send(m)
}
