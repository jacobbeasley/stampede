package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Todo is used by pop to map your todoes database table to your go code.
type Todo struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id" db:"organization_id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	IsCompleted    bool      `json:"is_completed" db:"is_completed"`
	Position       int       `json:"position" db:"position"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (t Todo) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Todoes is not required by pop and may be deleted
type Todoes []Todo

// String is not required by pop and may be deleted
func (t Todoes) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Todo) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.Title, Name: "Title"},
		&validators.UUIDIsPresent{Field: t.UserID, Name: "UserID"},
		&validators.UUIDIsPresent{Field: t.OrganizationID, Name: "OrganizationID"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (t *Todo) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (t *Todo) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
