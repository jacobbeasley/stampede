package actions

import (
	"fmt"

	"buffalo-app/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"golang.org/x/crypto/bcrypt"
)

// ProfileEdit renders the profile edit form.
func ProfileEdit(c buffalo.Context) error {
	user := c.Value("current_user").(*models.User)
	c.Set("user", user)
	return c.Render(200, r.HTML("profile/edit.plush.html"))
}

// ProfileUpdate processes the profile edit form.
func ProfileUpdate(c buffalo.Context) error {
	user := c.Value("current_user").(*models.User)
	tx := c.Value("tx").(*pop.Connection)

	// Bind to a temporary struct to prevent mass-assignment vulnerability
	formData := struct {
		FirstName   string `form:"FirstName"`
		LastName    string `form:"LastName"`
		Email       string `form:"Email"`
		PhoneNumber string `form:"PhoneNumber"`
	}{}

	if err := c.Bind(&formData); err != nil {
		return err
	}

	// Update only allowed fields
	if formData.FirstName != "" {
		user.FirstName.Valid = true
		user.FirstName.String = formData.FirstName
	} else {
		user.FirstName.Valid = false
		user.FirstName.String = ""
	}

	if formData.LastName != "" {
		user.LastName.Valid = true
		user.LastName.String = formData.LastName
	} else {
		user.LastName.Valid = false
		user.LastName.String = ""
	}

	if formData.PhoneNumber != "" {
		user.PhoneNumber.Valid = true
		user.PhoneNumber.String = formData.PhoneNumber
	} else {
		user.PhoneNumber.Valid = false
		user.PhoneNumber.String = ""
	}

	emailChanged := false
	if formData.Email != "" && formData.Email != user.Email {
		user.PendingEmail.Valid = true
		user.PendingEmail.String = formData.Email
		if err := user.GenerateEmailConfirmationToken(); err != nil {
			return err
		}
		emailChanged = true
	} else if formData.Email == "" {
		// Ensure empty email fails validation
		user.Email = ""
	} else {
		user.PendingEmail.Valid = false
		user.PendingEmail.String = ""
	}

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Set("user", user)
		c.Set("errors", verrs)
		c.Flash().Add("danger", "There was an error updating your profile.")
		return c.Render(422, r.HTML("profile/edit.plush.html"))
	}

	if emailChanged {
		c.Flash().Add("success", "Profile updated. Please check your old email to confirm your new email address.")
		// TODO: actually send email (we will add this later if mailer is set up, for now simulate the instruction)
	} else {
		c.Flash().Add("success", "Your profile was successfully updated.")
	}

	return c.Redirect(302, "/profile")
}

// ProfilePasswordUpdate processes the password update form.
func ProfilePasswordUpdate(c buffalo.Context) error {
	user := c.Value("current_user").(*models.User)
	tx := c.Value("tx").(*pop.Connection)

	// Ensure we have a fresh copy to avoid mutating the session user object if we fail
	u := &models.User{}
	if err := tx.Find(u, user.ID); err != nil {
		return err
	}

	currentPassword := c.Param("current_password")
	newPassword := c.Param("Password")
	passwordConfirmation := c.Param("PasswordConfirmation")

	// Validate current password
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(currentPassword)); err != nil {
		c.Flash().Add("danger", "Current password is incorrect.")
		return c.Redirect(302, "/profile")
	}

	if newPassword != passwordConfirmation {
		c.Flash().Add("danger", "New passwords do not match.")
		return c.Redirect(302, "/profile")
	}

	if newPassword == "" {
		c.Flash().Add("danger", "New password cannot be empty.")
		return c.Redirect(302, "/profile")
	}

	u.Password = newPassword
	u.PasswordConfirmation = passwordConfirmation

	if err := u.BeforeUpdate(tx); err != nil {
	    return err
	}

	verrs, err := tx.ValidateAndUpdate(u)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		c.Flash().Add("danger", fmt.Sprintf("Error updating password: %v", verrs.Error()))
		return c.Redirect(302, "/profile")
	}

	if err := tx.UpdateColumns(u, "password_hash"); err != nil {
		return err
	}

	c.Flash().Add("success", "Your password was successfully updated.")
	return c.Redirect(302, "/profile")
}
