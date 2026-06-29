package models

import (
	"encoding/json"
	"net/mail"
	"strings"
	"time"

	"crypto/rand"
	"encoding/hex"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v6"

	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User is used by pop to map your users database table to your go code.
type User struct {
	ID                         uuid.UUID     `json:"id" db:"id"`
	Email                      string        `json:"email" db:"email"`
	FirstName                  nulls.String  `json:"first_name" db:"first_name"`
	LastName                   nulls.String  `json:"last_name" db:"last_name"`
	PhoneNumber                nulls.String  `json:"phone_number" db:"phone_number"`
	PasswordHash               string        `json:"password_hash" db:"password_hash"`
	Password                   string        `json:"-" db:"-"`
	PasswordConfirmation       string        `json:"-" db:"-"`
	Roles                      []Role        `json:"roles,omitempty" many_to_many:"user_roles"`
	Organizations              Organizations `json:"organizations,omitempty" many_to_many:"user_roles"`
	ResetToken                 nulls.String  `json:"reset_token" db:"reset_token"`
	ResetTokenExpiresAt        nulls.Time    `json:"reset_token_expires_at" db:"reset_token_expires_at"`
	PendingEmail               nulls.String  `json:"pending_email" db:"pending_email"`
	EmailConfirmationToken     nulls.String  `json:"email_confirmation_token" db:"email_confirmation_token"`
	EmailConfirmationExpiresAt nulls.Time    `json:"email_confirmation_expires_at" db:"email_confirmation_expires_at"`
	FailedLoginAttempts        int           `json:"failed_login_attempts" db:"failed_login_attempts"`
	LastFailedLoginAt          nulls.Time    `json:"last_failed_login_at" db:"last_failed_login_at"`
	LastPasswordResetRequestAt nulls.Time    `json:"last_password_reset_request_at" db:"last_password_reset_request_at"`
	AccountVerified            bool          `json:"account_verified" db:"account_verified"`
	CreatedAt                  time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time     `json:"updated_at" db:"updated_at"`
}

// GenerateEmailConfirmationToken creates a random token and sets the expiration.
func (u *User) GenerateEmailConfirmationToken() error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	u.EmailConfirmationToken = nulls.NewString(hex.EncodeToString(bytes))
	u.EmailConfirmationExpiresAt = nulls.NewTime(time.Now().Add(1 * time.Hour))
	return nil
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// DisplayName returns first name + last name if either is present and non-empty, otherwise email.
func (u User) DisplayName() string {
	var parts []string
	if u.FirstName.Valid && strings.TrimSpace(u.FirstName.String) != "" {
		parts = append(parts, strings.TrimSpace(u.FirstName.String))
	}
	if u.LastName.Valid && strings.TrimSpace(u.LastName.String) != "" {
		parts = append(parts, strings.TrimSpace(u.LastName.String))
	}
	if len(parts) > 0 {
		return strings.Join(parts, " ")
	}
	return u.Email
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// GenerateResetToken creates a random reset token and sets the expiration.
func (u *User) GenerateResetToken() error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	u.ResetToken = nulls.NewString(hex.EncodeToString(bytes))
	u.ResetTokenExpiresAt = nulls.NewTime(time.Now().UTC().Add(72 * time.Hour))
	return nil
}

// Create validates and creates a new User.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	return tx.ValidateAndCreate(u)
}

func (u *User) BeforeCreate(tx *pop.Connection) error {
	return u.BeforeValidate(tx)
}

func (u *User) BeforeUpdate(tx *pop.Connection) error {
	return u.BeforeValidate(tx)
}

func (u *User) BeforeValidate(tx *pop.Connection) error {
	u.Email = strings.ToLower(u.Email)
	if u.Password != "" {
		ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(ph)
	}
	return nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "EmailFormat",
			Message: "%s is not a valid email address",
			Fn: func() bool {
				_, err := mail.ParseAddress(u.Email)
				return err == nil
			},
		},
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	// Generate reset token on create
	u.AccountVerified = false
	u.GenerateResetToken()
	// Password is no longer required on creation due to email verification flow
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	if u.Password != "" {
		var err error
		return validate.Validate(
			&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
		), err
	}
	return validate.NewErrors(), nil
}

// HasRole returns true if the user has the specified role ID.
func (u User) HasRole(roleID string) bool {
	for _, r := range u.Roles {
		if r.ID == roleID {
			return true
		}
	}
	return false
}

// IsSuperAdmin returns true if the user has the SUPER_ADMIN role.
func (u User) IsSuperAdmin() bool {
	return u.HasRole("SUPER_ADMIN")
}
