package actions

import (
	"buffalo-app/mailers"
	"github.com/gobuffalo/buffalo/worker"
)

var w worker.Worker

func init() {
	w = worker.NewSimple()

	w.Register("send_password_reset", func(args worker.Args) error {
		email := args["email"].(string)
		link := args["link"].(string)
		return mailers.SendPasswordReset(email, link)
	})

	w.Register("send_invitation", func(args worker.Args) error {
		email := args["email"].(string)
		organization := args["organization"].(string)
		return mailers.SendInvitation(email, organization)
	})

	w.Register("send_admin_password_reset", func(args worker.Args) error {
		email := args["email"].(string)
		link := args["link"].(string)
		return mailers.SendAdminPasswordReset(email, link)
	})
}
