package actions

import (
	"buffalo-app/models"
	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
)

// Helper to create and login a user
func (as *ActionSuite) loginUser() *models.User {
	user := &models.User{
		Email:                "testuser@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
		FirstName:            nulls.NewString("Test"),
		LastName:             nulls.NewString("User"),
		AccountVerified:      true,
	}
	as.NoError(as.DB.Create(user))

	// Create mock org to bypass Authorize check
	org := &models.Organization{ID: uuid.FromStringOrNil("00000000-0000-0000-0000-000000000001"), Name: "Mock Org"}
	as.DB.Create(org)

	// Mocking session
	as.Session.Set("current_user_id", user.ID)
	as.Session.Set("current_organization_id", org.ID)
	as.NoError(as.Session.Save())

	return user
}

// Add authenticity token to post requests
func (as *ActionSuite) postWithCSRF(url string, data map[string]string) {
	data["authenticity_token"] = "test-token"
}

func (as *ActionSuite) Test_ProfileEdit_RedirectsIfNotLoggedIn() {
	res := as.HTML("/profile").Get()
	as.Equal(303, res.Code)
	as.Contains(res.Location(), "/login")
}

func (as *ActionSuite) Test_ProfileEdit_LoadsWhenLoggedIn() {
	user := as.loginUser()

	res := as.HTML("/profile").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Profile Management")
	as.Contains(res.Body.String(), user.Email)
}

func (as *ActionSuite) Test_ProfileUpdate_Success() {
	user := as.loginUser()

	data := map[string]string{
		"FirstName":          "New",
		"LastName":           "Name",
		"PhoneNumber":        "555-555-5555",
		"Email":              user.Email,
		"authenticity_token": "test-token",
	}

	res := as.HTML("/profile").Post(data)

	as.Equal(302, res.Code)
	as.Equal("/profile", res.Location())

	// Verify in DB
	err := as.DB.Reload(user)
	as.NoError(err)
	as.True(user.FirstName.Valid)
	as.Equal("New", user.FirstName.String)
	as.True(user.LastName.Valid)
	as.Equal("Name", user.LastName.String)
	as.True(user.PhoneNumber.Valid)
	as.Equal("555-555-5555", user.PhoneNumber.String)
}

func (as *ActionSuite) Test_ProfileUpdate_EmailChangeInitiatesConfirmation() {
	user := as.loginUser()

	data := map[string]string{
		"FirstName":          "New",
		"LastName":           "Name",
		"Email":              "newemail@example.com",
		"authenticity_token": "test-token",
	}

	res := as.HTML("/profile").Post(data)

	as.Equal(302, res.Code)
	as.Equal("/profile", res.Location())

	// Verify in DB
	err := as.DB.Reload(user)
	as.NoError(err)
	as.Equal("testuser@example.com", user.Email) // Email should not change yet
	as.True(user.PendingEmail.Valid)
	as.Equal("newemail@example.com", user.PendingEmail.String)
	as.True(user.EmailConfirmationToken.Valid)
}

func (as *ActionSuite) Test_ProfileUpdate_FailureEmptyEmail() {
	as.loginUser()

	data := map[string]string{
		"FirstName":          "New",
		"LastName":           "Name",
		"Email":              "",
		"authenticity_token": "test-token",
	}

	res := as.HTML("/profile").Post(data)

	as.Equal(422, res.Code)
}

func (as *ActionSuite) Test_ProfilePasswordUpdate_Success() {
	user := as.loginUser()
	oldHash := user.PasswordHash

	data := map[string]string{
		"current_password":     "password123",
		"Password":             "newpassword123",
		"PasswordConfirmation": "newpassword123",
		"authenticity_token":   "test-token",
	}

	res := as.HTML("/profile/password").Post(data)

	as.Equal(302, res.Code)
	as.Equal("/profile", res.Location())

	// Verify password changed
	err := as.DB.Reload(user)
	as.NoError(err)
	as.NotEqual(oldHash, user.PasswordHash)
}

func (as *ActionSuite) Test_ProfilePasswordUpdate_WrongCurrentPassword() {
	as.loginUser()

	data := map[string]string{
		"current_password":     "wrongpassword",
		"Password":             "newpassword123",
		"PasswordConfirmation": "newpassword123",
		"authenticity_token":   "test-token",
	}

	res := as.HTML("/profile/password").Post(data)

	as.Equal(302, res.Code)
	as.Equal("/profile", res.Location())
}

func (as *ActionSuite) Test_ProfilePasswordUpdate_MismatchNewPassword() {
	as.loginUser()

	data := map[string]string{
		"current_password":     "password123",
		"Password":             "newpassword123",
		"PasswordConfirmation": "wrongconfirmation",
		"authenticity_token":   "test-token",
	}

	res := as.HTML("/profile/password").Post(data)

	as.Equal(302, res.Code)
	as.Equal("/profile", res.Location())
}
