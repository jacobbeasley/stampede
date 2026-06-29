package actions

import (
	"net/http"
	"sync"

	"buffalo-app/locales"
	"buffalo-app/models"
	"buffalo-app/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/v3/pop/popmw"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/csrf"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	appOnce.Do(func() {
		env := envy.Get("GO_ENV", "development")
		port := envy.Get("PORT", "3000")

		app = buffalo.New(buffalo.Options{
			Env:         env,
			SessionName: "_buffalo_test_session",
		})

		app.Logger.Infof("Booting application in %s mode on port %s", env, port)

		// Automatically redirect to SSL
		if envy.Get("FORCE_SSL", "false") == "true" {
			app.Use(forceSSL())
		}

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		if env != "test" {
			app.Use(csrf.New)
		} else {
			app.Use(func(next buffalo.Handler) buffalo.Handler {
				return func(c buffalo.Context) error {
					c.Set("authenticity_token", "test-token")
					return next(c)
				}
			})
		}

		// Wraps each request in a transaction, if enabled
		if envy.Get("DISABLE_DB_TRANSACTIONS", "false") != "true" {
			app.Use(popmw.Transaction(models.DB))
		}

		// Setup background worker
		app.Worker = w

		// Setup and use translations:
		app.Use(translations())
		app.Use(SetCurrentOrganization)
		app.Use(SetCurrentUser)

		setupRoutes(app)

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	})

	return app
}

func setupRoutes(app *buffalo.App) {
	// Health/Readiness endpoints (bypass auth, CSRF, etc.)
	apiBaseGroup := app.Group("/api")
	apiBaseGroup.Middleware.Skip(csrf.New, HealthCheck, ReadyCheck)
	if envy.Get("DISABLE_DB_TRANSACTIONS", "false") != "true" {
		apiBaseGroup.Middleware.Skip(popmw.Transaction(models.DB), HealthCheck, ReadyCheck)
	}
	apiBaseGroup.GET("/health", HealthCheck)
	apiBaseGroup.GET("/ready", ReadyCheck)

	app.GET("/", HomeHandler)

	app.GET("/register", AuthRegisterGet).Name("register")
	app.POST("/register", AuthRegisterPost).Name("register")

	app.GET("/login", AuthLoginGet)
	app.POST("/login", AuthLoginPost)
	app.GET("/auth/select_organization", Authorize(AuthSelectOrganizationGet)).Name("authSelectOrganization")
	app.POST("/auth/select_organization", Authorize(AuthSelectOrganizationPost)).Name("authSelectOrganization")
	app.GET("/logout", AuthLogout)

	app.GET("/password/reset", PasswordResetGet)
	app.POST("/password/reset", PasswordResetPost)
	app.GET("/password/edit/{user_id}", PasswordEditGet)
	app.POST("/password/edit/{user_id}", PasswordEditPost)

	app.GET("/profile", Authorize(ProfileEdit)).Name("profileEdit")
	app.POST("/profile", Authorize(ProfileUpdate)).Name("profileUpdate")
	app.POST("/profile/password", Authorize(ProfilePasswordUpdate)).Name("profilePasswordUpdate")

	app.GET("/todos", Authorize(TodosIndex))

	apiGroup := app.Group("/api/todos")
	apiGroup.Use(Authorize)
	apiGroup.GET("/", TodosList)
	apiGroup.POST("/", TodosCreate)
	apiGroup.PUT("/reorder", TodosReorder)
	apiGroup.PUT("/{todo_id}", TodosUpdate)
	apiGroup.DELETE("/{todo_id}", TodosDelete)

	adminGroup := app.Group("/admin")
	adminGroup.Use(Authorize)
	adminGroup.Use(AuthorizeAdmin)
	adminUsersResource := AdminUsersResource{}
	adminGroup.Resource("/users", adminUsersResource)
	adminGroup.POST("/users/{user_id}/password_reset", adminUsersResource.PasswordReset)

	superAdminGroup := adminGroup.Group("/super")
	superAdminGroup.Use(AuthorizeSuperAdmin)
	superAdminGroup.Resource("/organizations", OrganizationsResource{})
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
