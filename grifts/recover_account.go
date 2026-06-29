package grifts

import (
	"fmt"
	"strings"

	"buffalo-app/models"

	"github.com/gobuffalo/grift/grift"
	"github.com/gobuffalo/nulls"
	"golang.org/x/crypto/bcrypt"
)

var _ = grift.Add("recover_account", func(c *grift.Context) error {
	if len(c.Args) < 2 {
		return fmt.Errorf("usage: buffalo task recover_account <email> <new_password>")
	}

	email := strings.ToLower(strings.TrimSpace(c.Args[0]))
	newPassword := c.Args[1]

	u := &models.User{}
	err := models.DB.Where("email = ?", email).First(u)
	if err != nil {
		return fmt.Errorf("user with email %s not found: %v", email, err)
	}

	ph, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate password hash: %v", err)
	}

	u.PasswordHash = string(ph)
	u.AccountVerified = true
	u.FailedLoginAttempts = 0
	u.LastFailedLoginAt = nulls.Time{}

	err = models.DB.Update(u)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	fmt.Printf("Account recovered successfully. User %s is now verified and password has been reset.\n", email)
	return nil
})
