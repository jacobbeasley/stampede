package actions

import (
	"os"
	"strings"

	"buffalo-app/public"
	"buffalo-app/templates"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/helpers/forms"
)

var r *render.Engine

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
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// fs.FS containing templates
		TemplatesFS: templates.FS(),

		// fs.FS containing assets
		AssetsFS: public.FS(),

		// Add template helpers here:
		Helpers: render.Helpers{
			"viteAsset": viteAsset,
			"viteCSS":   viteCSS,
			"form_for":  forms.FormFor,
			"formFor":   forms.FormFor,
			"GO_ENV":    envy.Get("GO_ENV", "development"),
			"siteURL":   siteURL,
		},
	})
}
