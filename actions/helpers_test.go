package actions

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SiteURL(t *testing.T) {
	// Keep existing environment clean
	origSiteURL := os.Getenv("SITE_URL")
	defer func() {
		if origSiteURL != "" {
			os.Setenv("SITE_URL", origSiteURL)
		} else {
			os.Unsetenv("SITE_URL")
		}
	}()

	// 1. Test default value
	os.Unsetenv("SITE_URL")
	assert.Equal(t, "http://localhost:3000/", siteURL())

	// 2. Test custom value with trailing slash
	os.Setenv("SITE_URL", "https://todoflow.example.com/")
	assert.Equal(t, "https://todoflow.example.com/", siteURL())

	// 3. Test custom value without trailing slash (should append it)
	os.Setenv("SITE_URL", "https://todoflow.example.com")
	assert.Equal(t, "https://todoflow.example.com/", siteURL())
}
