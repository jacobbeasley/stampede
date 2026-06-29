package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"buffalo-app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthRegisterGet renders the registration form
func AuthRegisterGet(c buffalo.Context) error {
	c.Set("user", &models.User{})
	return c.Render(http.StatusOK, r.HTML("auth/register.plush.html"))
}

// AuthRegisterPost handles the registration form submission
func AuthRegisterPost(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(u)
	if err != nil {
		return err
	}

	u.AccountVerified = false
	if err := u.GenerateResetToken(); err != nil {
		return err
	}

	if err := tx.Update(u); err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(http.StatusOK, r.HTML("auth/register.plush.html"))
	}

	// Default Organization ID
	defaultOrgID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")

	// Assign 'USER' role to newly registered users in Default Organization
	userRole := &models.UserRole{
		UserID:         u.ID,
		RoleID:         "USER",
		OrganizationID: defaultOrgID,
	}
	if err := tx.Create(userRole); err != nil {
		return err
	}

	// Send verification email asynchronously
	link := fmt.Sprintf("%spassword/edit/%s?token=%s", siteURL(), u.ID, u.ResetToken.String)
	w.Perform(worker.Job{
		Queue:   "default",
		Handler: "send_verification",
		Args: worker.Args{
			"email": u.Email,
			"link":  link,
		},
	})

	c.Flash().Add("success", "Account created! Please check your email to verify your account and set a password.")

	return c.Redirect(http.StatusSeeOther, "/login")
}

// AuthLoginGet renders the login form
func AuthLoginGet(c buffalo.Context) error {
	c.Set("user", &models.User{})
	return c.Render(http.StatusOK, r.HTML("auth/login.plush.html"))
}

// AuthLoginPost handles the login form submission
func AuthLoginPost(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	password := u.Password
	err := tx.Where("email = ?", strings.ToLower(u.Email)).First(u)
	if err != nil {
		if err == sql.ErrNoRows {
			c.Set("user", u)
			c.Flash().Add("danger", "Invalid email or password.")
			return c.Render(http.StatusOK, r.HTML("auth/login.plush.html"))
		}
		return err
	}

	if !u.AccountVerified {
		if err := u.GenerateResetToken(); err != nil {
			return err
		}
		if err := tx.Update(u); err != nil {
			return err
		}

		link := fmt.Sprintf("%spassword/edit/%s?token=%s", siteURL(), u.ID, u.ResetToken.String)
		w.Perform(worker.Job{
			Queue:   "default",
			Handler: "send_verification",
			Args: worker.Args{
				"email": u.Email,
				"link":  link,
			},
		})

		c.Flash().Add("danger", "Your account is unverified. A new verification email has been sent.")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if u.FailedLoginAttempts >= 5 && u.LastFailedLoginAt.Valid {
		lockoutExpiration := u.LastFailedLoginAt.Time.Add(5 * time.Minute)
		if time.Now().Before(lockoutExpiration) {
			c.Flash().Add("danger", "Account locked for 5 minutes due to too many failed login attempts.")
			return c.Redirect(http.StatusSeeOther, "/login")
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		u.FailedLoginAttempts++
		u.LastFailedLoginAt = nulls.NewTime(time.Now())
		tx.UpdateColumns(u, "failed_login_attempts", "last_failed_login_at")

		c.Set("user", u)
		c.Flash().Add("danger", "Invalid email or password.")
		return c.Render(http.StatusOK, r.HTML("auth/login.plush.html"))
	}

	u.FailedLoginAttempts = 0
	u.LastFailedLoginAt = nulls.Time{}
	tx.UpdateColumns(u, "failed_login_attempts", "last_failed_login_at")

	rotateSession(c)
	c.Session().Set("current_user_id", u.ID)

	// Load user's organizations
	err = tx.Load(u, "Organizations")
	if err != nil {
		return err
	}

	if len(u.Organizations) == 0 {
		// This shouldn't happen if user is correctly created, but handle it
		c.Flash().Add("danger", "You are not a member of any organization.")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if len(u.Organizations) == 1 {
		c.Session().Set("current_organization_id", u.Organizations[0].ID)
		c.Flash().Add("success", "Welcome back to TodoFlow!")
		return c.Redirect(http.StatusSeeOther, "/todos")
	}

	// Multiple organizations, redirect to selection
	return c.Redirect(http.StatusSeeOther, "/auth/select_organization")
}

// AuthSelectOrganizationGet renders the organization selection form
func AuthSelectOrganizationGet(c buffalo.Context) error {
	u := c.Value("current_user").(*models.User)
	tx := c.Value("tx").(*pop.Connection)
	err := tx.Load(u, "Organizations")
	if err != nil {
		return err
	}
	c.Set("organizations", u.Organizations)
	return c.Render(http.StatusOK, r.HTML("auth/select_organization.plush.html"))
}

// AuthSelectOrganizationPost handles the organization selection
func AuthSelectOrganizationPost(c buffalo.Context) error {
	orgID := c.Param("organization_id")
	if orgID == "" {
		return c.Redirect(http.StatusSeeOther, "/auth/select_organization")
	}

	// Verify user is member of the selected org
	u := c.Value("current_user").(*models.User)
	tx := c.Value("tx").(*pop.Connection)
	exists, err := tx.Where("user_id = ? AND organization_id = ?", u.ID, orgID).Exists(&models.UserRole{})
	if err != nil {
		return err
	}

	isSuperAdmin := false
	for _, r := range u.Roles {
		if r.ID == "SUPER_ADMIN" {
			isSuperAdmin = true
			break
		}
	}

	if !exists && !isSuperAdmin {
		c.Flash().Add("danger", "You are not a member of that organization.")
		return c.Redirect(http.StatusSeeOther, "/auth/select_organization")
	}

	c.Session().Set("current_organization_id", orgID)
	c.Flash().Add("success", "Organization selected.")

	redirectURL := c.Param("redirect_url")
	if redirectURL == "" {
		redirectURL = "/todos"
	}
	return c.Redirect(http.StatusSeeOther, redirectURL)
}

func rotateSession(c buffalo.Context) {
	s := c.Session()
	s.Session.Options.MaxAge = -1
	s.Session.Save(c.Request(), c.Response())
	s.Session.ID = ""
	s.Session.IsNew = true
	s.Session.Options.MaxAge = 0
	s.Session.Values = make(map[interface{}]interface{})
}

// AuthLogout clears the session
func AuthLogout(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out.")
	return c.Redirect(http.StatusSeeOther, "/")
}

// SetCurrentUser attempts to find a user from the session and bind it to context
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)

			// Load user and their roles for the current organization if set
			orgID := c.Session().Get("current_organization_id")

			var err error
			if orgID != nil {
				err = tx.Eager("Roles").Where("id = ?", uid).First(u)
				if err == nil {
					// Manually load roles filtered by organization
					roles := []models.Role{}
					err = tx.RawQuery("SELECT r.* FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = ? AND ur.organization_id = ?", uid, orgID).All(&roles)
					if err == nil {
						u.Roles = roles
					}

					// Ensure SUPER_ADMIN role is retained globally if they have it in any organization
					superAdminRoles := []models.Role{}
					err = tx.RawQuery("SELECT DISTINCT r.* FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = ? AND r.id = 'SUPER_ADMIN'", uid).All(&superAdminRoles)
					if err == nil && len(superAdminRoles) > 0 {
						hasSuperAdmin := false
						for _, r := range u.Roles {
							if r.ID == "SUPER_ADMIN" {
								hasSuperAdmin = true
								break
							}
						}
						if !hasSuperAdmin {
							u.Roles = append(u.Roles, superAdminRoles[0])
						}
					}
				}
			} else {
				err = tx.Eager("Roles").Find(u, uid)
			}

			if err != nil {
				c.Logger().Errorf("SetCurrentUser error loading user (uid=%v, orgID=%v): %v", uid, orgID, err)
			} else {
				c.Set("current_user", u)

				// Also load all user's organizations for the switcher
				tx.Load(u, "Organizations")
				c.Set("user_organizations", u.Organizations)
			}
		}
		return next(c)
	}
}

// SetCurrentOrganization attempts to find an organization from the session and bind it to context
func SetCurrentOrganization(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if oid := c.Session().Get("current_organization_id"); oid != nil {
			o := &models.Organization{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(o, oid)
			if err == nil {
				c.Set("current_organization", o)
			}
		}
		return next(c)
	}
}

// Authorize requires a user to be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if c.Value("current_user") == nil {
			c.Session().Clear()
			c.Flash().Add("danger", "You must be logged in to view that page.")
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		// Ensure organization is selected for non-auth routes
		path := strings.TrimRight(c.Request().URL.Path, "/")
		if path != "/auth/select_organization" && path != "/logout" && path != "/register" && path != "/login" && c.Session().Get("current_organization_id") == nil {
			return c.Redirect(http.StatusSeeOther, "/auth/select_organization")
		}

		return next(c)
	}
}

// AuthorizeAdmin requires a user to be an admin to access a route
func AuthorizeAdmin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if u, ok := c.Value("current_user").(*models.User); ok {
			for _, role := range u.Roles {
				if role.ID == "ADMIN" || role.ID == "SUPER_ADMIN" {
					return next(c)
				}
			}
		}
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(http.StatusSeeOther, "/")
	}
}

// AuthorizeSuperAdmin requires a user to be a super admin to access a route
func AuthorizeSuperAdmin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if u, ok := c.Value("current_user").(*models.User); ok {
			for _, role := range u.Roles {
				if role.ID == "SUPER_ADMIN" {
					return next(c)
				}
			}
		}
		c.Flash().Add("danger", "You are not authorized to view that page.")
		return c.Redirect(http.StatusSeeOther, "/")
	}
}
