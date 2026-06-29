package actions

import (
	"buffalo-app/models"

	"github.com/gobuffalo/buffalo"
)

// ReadyCheck immediately returns a 200 OK status to indicate the app process is running.
// This is typically used as a liveness probe.
func ReadyCheck(c buffalo.Context) error {
	return c.Render(200, r.JSON(map[string]string{
		"status": "ready",
	}))
}

// HealthCheck verifies connectivity to critical infrastructure components (like the DB).
// This is typically used as a readiness probe.
func HealthCheck(c buffalo.Context) error {
	if err := models.DB.RawQuery("SELECT 1").Exec(); err != nil {
		c.Logger().Errorf("Health check failed: database connection error: %v", err)
		return c.Render(503, r.JSON(map[string]string{
			"status": "error",
			"error":  "database unavailable",
		}))
	}

	return c.Render(200, r.JSON(map[string]string{
		"status": "healthy",
	}))
}
