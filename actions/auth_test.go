package actions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"buffalo-app/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/stretchr/testify/require"
)

func Test_AuthorizeAdmin_Middleware(t *testing.T) {
	ren := render.New(render.Options{})

	// Dummy handler that just returns success
	handler := func(c buffalo.Context) error {
		return c.Render(http.StatusOK, ren.String("success"))
	}

	// Middleware to set user
	setUserMiddleware := func(user *models.User) buffalo.MiddlewareFunc {
		return func(next buffalo.Handler) buffalo.Handler {
			return func(c buffalo.Context) error {
				if user != nil {
					c.Set("current_user", user)
				}
				return next(c)
			}
		}
	}

	t.Run("No User is Redirected", func(t *testing.T) {
		req := require.New(t)
		app := buffalo.New(buffalo.Options{SessionName: "_test_session"})
		app.Use(setUserMiddleware(nil))
		app.Use(AuthorizeAdmin)
		app.GET("/admin-only", handler)

		r := httptest.NewRequest("GET", "/admin-only", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)

		req.Equal(http.StatusSeeOther, w.Code)
		req.Equal("/", w.Header().Get("Location"))
	})

	t.Run("Non-Admin User is Redirected", func(t *testing.T) {
		req := require.New(t)
		app := buffalo.New(buffalo.Options{SessionName: "_test_session"})
		app.Use(setUserMiddleware(&models.User{Roles: []models.Role{{ID: "USER"}}}))
		app.Use(AuthorizeAdmin)
		app.GET("/admin-only", handler)

		r := httptest.NewRequest("GET", "/admin-only", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)

		req.Equal(http.StatusSeeOther, w.Code)
		req.Equal("/", w.Header().Get("Location"))
	})

	t.Run("Admin User is Allowed", func(t *testing.T) {
		req := require.New(t)
		app := buffalo.New(buffalo.Options{SessionName: "_test_session"})
		app.Use(setUserMiddleware(&models.User{Roles: []models.Role{{ID: "ADMIN"}}}))
		app.Use(AuthorizeAdmin)
		app.GET("/admin-only", handler)

		r := httptest.NewRequest("GET", "/admin-only", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, r)

		req.Equal(http.StatusOK, w.Code)
	})
}
