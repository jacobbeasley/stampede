package models

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/gobuffalo/nulls"
)

func (ms *ModelSuite) Test_User_BeforeCreate() {
	u := &User{
		Email:    "TEST@EXAMPLE.COM",
		Password: "password123",
	}

	err := u.BeforeCreate(ms.DB)
	ms.NoError(err)
	ms.Equal("test@example.com", u.Email)
	ms.NotEmpty(u.PasswordHash)

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("password123"))
	ms.NoError(err)
}

func (ms *ModelSuite) Test_User_BeforeCreate_NoPassword() {
	u := &User{
		Email: "TEST2@EXAMPLE.COM",
	}

	err := u.BeforeCreate(ms.DB)
	ms.NoError(err)
	ms.Equal("test2@example.com", u.Email)
	ms.Empty(u.PasswordHash)
}

func (ms *ModelSuite) Test_User_DisplayName() {
	// Import nulls from github.com/gobuffalo/nulls if not done already, but user.go has it.
	// Since user_test.go is in the models package, we can access nulls.String directly.
	u := User{
		Email: "test@example.com",
	}
	ms.Equal("test@example.com", u.DisplayName())

	u.FirstName = nulls.NewString("John")
	ms.Equal("John", u.DisplayName())

	u.LastName = nulls.NewString("Doe")
	ms.Equal("John Doe", u.DisplayName())

	u.FirstName = nulls.NewString("")
	ms.Equal("Doe", u.DisplayName())

	u.FirstName = nulls.NewString(" ")
	ms.Equal("Doe", u.DisplayName())
}
