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
)

// PasswordResetGet renders the form to request a password reset
func PasswordResetGet(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("auth/password_reset.plush.html"))
}

// PasswordResetPost handles the request to send a password reset link
func PasswordResetPost(c buffalo.Context) error {
	email := strings.ToLower(strings.TrimSpace(c.Param("email")))
	if email == "" {
		c.Flash().Add("danger", "Email is required.")
		return c.Render(http.StatusOK, r.HTML("auth/password_reset.plush.html"))
	}

	tx := c.Value("tx").(*pop.Connection)
	u := &models.User{}
	err := tx.Where("email = ?", email).First(u)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't leak user existence
			c.Flash().Add("success", "If an account exists with that email, a password reset link has been sent. Note that password reset requests are rate limited to once every 5 minutes.")
			return c.Redirect(http.StatusSeeOther, "/login")
		}
		return err
	}

	if u.LastPasswordResetRequestAt.Valid && time.Since(u.LastPasswordResetRequestAt.Time) < 5*time.Minute {
		// Rate limit reached, silently ignore to prevent spam but show same success message
		c.Flash().Add("success", "If an account exists with that email, a password reset link has been sent. Note that password reset requests are rate limited to once every 5 minutes.")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if err := u.GenerateResetToken(); err != nil {
		return err
	}

	u.LastPasswordResetRequestAt = nulls.NewTime(time.Now())

	if err := tx.Update(u); err != nil {
		return err
	}

	// Send email asynchronously
	link := fmt.Sprintf("%spassword/edit/%s?token=%s", siteURL(), u.ID, u.ResetToken.String)
	w.Perform(worker.Job{
		Queue:   "default",
		Handler: "send_password_reset",
		Args: worker.Args{
			"email": u.Email,
			"link":  link,
		},
	})

	c.Flash().Add("success", "If an account exists with that email, a password reset link has been sent. Note that password reset requests are rate limited to once every 5 minutes.")
	return c.Redirect(http.StatusSeeOther, "/login")
}

func findUserByResetToken(c buffalo.Context) (*models.User, string, error) {
	userID := c.Param("user_id")
	token := c.Param("token")
	if token == "" || userID == "" {
		c.Flash().Add("danger", "Invalid or missing token.")
		return nil, token, c.Redirect(http.StatusSeeOther, "/login")
	}

	tx := c.Value("tx").(*pop.Connection)
	u := &models.User{}
	err := tx.Find(u, userID)
	if err != nil {
		c.Flash().Add("danger", "Invalid or expired token.")
		return nil, token, c.Redirect(http.StatusSeeOther, "/login")
	}

	if !u.ResetToken.Valid || u.ResetToken.String != token {
		// Wrong token submitted for this user, clear the reset token to prevent brute force
		u.ResetToken = nulls.String{}
		u.ResetTokenExpiresAt = nulls.Time{}
		tx.Update(u)

		c.Flash().Add("danger", "Invalid or expired token.")
		return nil, token, c.Redirect(http.StatusSeeOther, "/login")
	}

	if u.ResetTokenExpiresAt.Valid && u.ResetTokenExpiresAt.Time.UTC().Before(time.Now().UTC()) {
		c.Flash().Add("danger", "Token has expired. Please request a new one.")
		return nil, token, c.Redirect(http.StatusSeeOther, "/password/reset")
	}

	return u, token, nil
}

// PasswordEditGet renders the form to set a new password, given a valid token
func PasswordEditGet(c buffalo.Context) error {
	u, token, err := findUserByResetToken(c)
	if err != nil {
		return err
	}
	if u == nil {
		return nil // Response was already handled (Redirected)
	}

	c.Set("token", token)
	c.Set("user", u)
	return c.Render(http.StatusOK, r.HTML("auth/password_edit.plush.html"))
}

// PasswordEditPost handles the submission of a new password
func PasswordEditPost(c buffalo.Context) error {
	u, token, err := findUserByResetToken(c)
	if err != nil {
		return err
	}
	if u == nil {
		return nil // Response was already handled (Redirected)
	}

	if err := c.Bind(u); err != nil {
		return err
	}

	u.ResetToken = nulls.String{}
	u.ResetTokenExpiresAt = nulls.Time{}
	u.AccountVerified = true
	u.FailedLoginAttempts = 0
	u.LastFailedLoginAt = nulls.Time{}

	tx := c.Value("tx").(*pop.Connection)

	// Ensure BeforeValidate runs
	if err := u.BeforeUpdate(tx); err != nil {
	    return err
	}

	verrs, err := tx.ValidateAndUpdate(u)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("token", token)
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(http.StatusOK, r.HTML("auth/password_edit.plush.html"))
	}

	// Save the fields explicitly since we're updating a limited set
	if err := tx.UpdateColumns(u, "reset_token", "reset_token_expires_at", "account_verified", "failed_login_attempts", "last_failed_login_at", "password_hash"); err != nil {
		return err
	}

	rotateSession(c)
	c.Session().Set("current_user_id", u.ID)

	// Load user's organizations
	err = tx.Load(u, "Organizations")
	if err != nil {
		return err
	}

	if len(u.Organizations) == 0 {
		c.Flash().Add("danger", "You are not a member of any organization.")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	c.Flash().Add("success", "Password successfully updated! You are now logged in.")

	if len(u.Organizations) == 1 {
		c.Session().Set("current_organization_id", u.Organizations[0].ID)
		return c.Redirect(http.StatusSeeOther, "/todos")
	}

	// Multiple organizations, redirect to selection
	return c.Redirect(http.StatusSeeOther, "/auth/select_organization")
}
