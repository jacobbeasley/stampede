package mailers

import (
	"log"
	"os"
	"strings"

	"buffalo-app/templates"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
)

var (
	smtp mail.Sender
	r    *render.Engine
)

// siteURL returns the configured SITE_URL env variable or defaults to "http://localhost:3000/".
// It ensures that the returned URL has a trailing slash.
func siteURL() string {
	url := os.Getenv("SITE_URL")
	if url == "" {
		url = "http://localhost:3000/"
	}
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return url
}

func init() {
	// Pulling config from the env.
	port := envy.Get("SMTP_PORT", "")
	host := envy.Get("SMTP_HOST", "")
	user := envy.Get("SMTP_USER", "")
	password := envy.Get("SMTP_PASSWORD", "")

	if host != "" && port != "" {
		var err error
		smtp, err = mail.NewSMTPSender(host, port, user, password)

		if err != nil {
			log.Fatal(err)
		}
	} else {
		// Fallback to logging the email if SMTP isn't configured
		smtp = LoggerSender{}
	}

	r = render.New(render.Options{
		HTMLLayout:  "mail/layout.plush.html",
		TemplatesFS: templates.FS(),
		Helpers: render.Helpers{
			"siteURL": siteURL,
		},
	})
}
