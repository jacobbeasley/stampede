package actions

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"buffalo-app/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/worker"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

type AdminUsersResource struct {
	buffalo.Resource
}

// findUser is a helper to find a user by admin_user_id parameter.
func (v AdminUsersResource) findUser(c buffalo.Context) (*models.User, error) {
	tx := c.Value("tx").(*pop.Connection)
	user := &models.User{}
	// Buffalo resource for AdminUsersResource uses admin_user_id
	id := c.Param("admin_user_id")
	if id == "" {
		id = c.Param("user_id") // fallback for custom routes
	}
	if err := tx.Eager("Roles").Find(user, id); err != nil {
		return nil, c.Error(http.StatusNotFound, err)
	}

	currentUser := c.Value("current_user").(*models.User)
	isSuperAdmin := currentUser.HasRole("SUPER_ADMIN")

	// Protect SUPER_ADMIN users from modification by tenant-level administrators
	if !isSuperAdmin {
		isTargetSuperAdmin := user.HasRole("SUPER_ADMIN")
		if isTargetSuperAdmin {
			return nil, c.Error(http.StatusForbidden, fmt.Errorf("you do not have permission to access a Super Administrator"))
		}

		// Enforce cross-tenant isolation
		orgID := c.Session().Get("current_organization_id")
		if orgID == nil {
			return nil, c.Error(http.StatusForbidden, fmt.Errorf("active organization not found in session"))
		}
		exists, err := tx.Where("user_id = ? AND organization_id = ?", user.ID, orgID).Exists(&models.UserRole{})
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, c.Error(http.StatusForbidden, fmt.Errorf("you do not have permission to access this user"))
		}
	}

	return user, nil
}

// List default implementation.
func (v AdminUsersResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	orgID := c.Session().Get("current_organization_id")
	if orgID == nil {
		return c.Redirect(http.StatusSeeOther, "/auth/select_organization")
	}

	users := &models.Users{}
	// Only show users belonging to the current organization
	q := tx.PaginateFromParams(c.Params())
	err := q.RawQuery("SELECT DISTINCT u.* FROM users u JOIN user_roles ur ON u.id = ur.user_id WHERE ur.organization_id = ?", orgID).All(users)
	if err != nil {
		return err
	}

	if len(*users) > 0 {
		userIDs := make([]uuid.UUID, len(*users))
		for i, u := range *users {
			userIDs[i] = u.ID
		}

		type userRoleWithInfo struct {
			models.Role
			UserID uuid.UUID `db:"user_id"`
		}
		var rolesWithInfo []userRoleWithInfo
		err = tx.RawQuery("SELECT r.*, ur.user_id FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.organization_id = ? AND ur.user_id IN (?)", orgID, userIDs).All(&rolesWithInfo)
		if err == nil {
			roleMap := make(map[uuid.UUID][]models.Role)
			for _, r := range rolesWithInfo {
				roleMap[r.UserID] = append(roleMap[r.UserID], r.Role)
			}
			for i := range *users {
				(*users)[i].Roles = roleMap[(*users)[i].ID]
			}
		}
	}

	c.Set("pagination", q.Paginator)
	c.Set("users", users)

	return c.Render(http.StatusOK, r.HTML("admin/users/index.plush.html"))
}

// Show default implementation.
func (v AdminUsersResource) Show(c buffalo.Context) error {
	user, err := v.findUser(c)
	if err != nil {
		return err
	}

	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("admin/users/show.plush.html"))
}

// New default implementation.
func (v AdminUsersResource) New(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	roles := &models.Roles{}
	if err := tx.All(roles); err != nil {
		return err
	}
	c.Set("roles", roles)
	c.Set("user", &models.User{})
	return c.Render(http.StatusOK, r.HTML("admin/users/new.plush.html"))
}

// Create default implementation.
func (v AdminUsersResource) Create(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	orgIDRaw := c.Session().Get("current_organization_id")
	if orgIDRaw == nil {
		return c.Redirect(http.StatusSeeOther, "/auth/select_organization")
	}
	orgID, err := uuid.FromString(fmt.Sprintf("%v", orgIDRaw))
	if err != nil {
		return err
	}

	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return err
	}

	// Sanitize empty string to nulls.String properly
	if user.FirstName.String == "" {
		user.FirstName.Valid = false
	} else {
		user.FirstName.Valid = true
	}
	if user.LastName.String == "" {
		user.LastName.Valid = false
	} else {
		user.LastName.Valid = true
	}

	// Check if user already exists
	existingUser := &models.User{}
	err = tx.Where("email = ?", user.Email).First(existingUser)
	if err != nil {
		// User doesn't exist, create them
		// Generate random password if none is provided, as we are creating a user without registration
		if user.Password == "" {
			b := make([]byte, 16)
			if _, err := rand.Read(b); err == nil {
				user.Password = hex.EncodeToString(b)
				user.PasswordConfirmation = user.Password
			}
		}

		verrs, err := user.Create(tx)
		if err != nil {
			return err
		}

		if verrs.HasAny() {
			roles := &models.Roles{}
			tx.All(roles)
			c.Set("roles", roles)
			c.Set("errors", verrs)
			c.Set("user", user)
			return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/users/new.plush.html"))
		}
	} else {
		// User exists, check if they are already in this organization
		exists, err := tx.Where("user_id = ? AND organization_id = ?", existingUser.ID, orgID).Exists(&models.UserRole{})
		if err != nil {
			return err
		}
		if exists {
			c.Flash().Add("danger", "User is already a member of this organization.")
			return c.Redirect(http.StatusSeeOther, "/admin/users")
		}
		user = existingUser
		w.Perform(worker.Job{
			Queue:   "default",
			Handler: "send_invitation",
			Args: worker.Args{
				"email":        user.Email,
				"organization": c.Value("current_organization").(*models.Organization).Name,
			},
		})
	}

	c.Request().ParseForm()
	roleIDs := c.Request().Form["RoleIDs"]
	if len(roleIDs) == 0 {
		// Try singular
		if r := c.Request().Form.Get("RoleIDs"); r != "" {
			roleIDs = []string{r}
		}
	}
	if len(roleIDs) == 0 {
		roleIDs = c.Request().Form["Roles"]
	}
	if len(roleIDs) == 0 {
		if r := c.Request().Form.Get("Roles"); r != "" {
			roleIDs = []string{r}
		}
	}
	currentUser := c.Value("current_user").(*models.User)
	isSuperAdmin := currentUser.HasRole("SUPER_ADMIN")

	for _, roleID := range roleIDs {
		// Security check: Admins cannot assign SUPER_ADMIN role
		if roleID == "SUPER_ADMIN" && !isSuperAdmin {
			continue
		}
		userRole := &models.UserRole{
			UserID:         user.ID,
			RoleID:         roleID,
			OrganizationID: orgID,
		}
		if err := tx.Create(userRole); err != nil {
			return err
		}
	}

	c.Flash().Add("success", "User was added successfully.")
	return c.Redirect(http.StatusSeeOther, "/admin/users")
}

// Edit default implementation.
func (v AdminUsersResource) Edit(c buffalo.Context) error {
	user, err := v.findUser(c)
	if err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	roles := &models.Roles{}
	if err := tx.All(roles); err != nil {
		return err
	}
	c.Set("roles", roles)
	c.Set("user", user)
	return c.Render(http.StatusOK, r.HTML("admin/users/edit.plush.html"))
}

// Update default implementation.
func (v AdminUsersResource) Update(c buffalo.Context) error {
	user, err := v.findUser(c)
	if err != nil {
		return err
	}

	// Bind to a temporary input struct to prevent direct mass-assignment onto the database model
	input := &models.User{}
	if err := c.Bind(input); err != nil {
		return err
	}

	// Copy only individual permitted fields to the searched database user model
	user.FirstName = input.FirstName
	if user.FirstName.String == "" {
		user.FirstName.Valid = false
	} else {
		user.FirstName.Valid = true
	}

	user.LastName = input.LastName
	if user.LastName.String == "" {
		user.LastName.Valid = false
	} else {
		user.LastName.Valid = true
	}

	user.Email = input.Email
	user.PhoneNumber = input.PhoneNumber

	currentUser := c.Value("current_user").(*models.User)
	isSelf := user.ID == currentUser.ID
	isSuperAdmin := currentUser.IsSuperAdmin()

	// Only Super Administrators or users editing their own profiles are allowed to modify passwords
	if input.Password != "" && (isSelf || isSuperAdmin) {
		user.Password = input.Password
		user.PasswordConfirmation = input.PasswordConfirmation
	}

	tx := c.Value("tx").(*pop.Connection)
	orgID := c.Session().Get("current_organization_id")

	verrs, err := tx.ValidateAndUpdate(user)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		roles := &models.Roles{}
		tx.All(roles)
		c.Set("roles", roles)
		c.Set("errors", verrs)
		c.Set("user", user)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/users/edit.plush.html"))
	}

	// Update user roles for this organization
	c.Request().ParseForm()
	newRoleIDs := c.Request().Form["RoleIDs"]
	if len(newRoleIDs) == 0 {
		// Try singular
		if r := c.Request().Form.Get("RoleIDs"); r != "" {
			newRoleIDs = []string{r}
		}
	}
	if len(newRoleIDs) == 0 {
		newRoleIDs = c.Request().Form["Roles"]
	}
	if len(newRoleIDs) == 0 {
		if r := c.Request().Form.Get("Roles"); r != "" {
			newRoleIDs = []string{r}
		}
	}

	isSuperAdmin = currentUser.HasRole("SUPER_ADMIN")
	isAdmin := currentUser.HasRole("ADMIN")

	// Remove existing roles for THIS organization only
	if err := tx.RawQuery("DELETE FROM user_roles WHERE user_id = ? AND organization_id = ?", user.ID, orgID).Exec(); err != nil {
		return err
	}

	orgIDUUID, err := uuid.FromString(fmt.Sprintf("%v", orgID))
	if err != nil {
		return err
	}

	// Add new roles
	for _, roleID := range newRoleIDs {
		// Security checks
		if roleID == "SUPER_ADMIN" && !isSuperAdmin {
			continue // Non-superadmins cannot grant SUPER_ADMIN
		}

		userRole := &models.UserRole{
			UserID:         user.ID,
			RoleID:         roleID,
			OrganizationID: orgIDUUID,
		}
		if err := tx.Create(userRole); err != nil {
			return err
		}
	}

	// Re-verify: Admin cannot remove their own admin role
	if isSelf && isAdmin {
		exists, _ := tx.Where("user_id = ? AND organization_id = ? AND role_id = 'ADMIN'", currentUser.ID, orgID).Exists(&models.UserRole{})
		if !exists {
			// Restore it
			tx.Create(&models.UserRole{UserID: currentUser.ID, OrganizationID: orgIDUUID, RoleID: "ADMIN"})
			c.Flash().Add("warning", "You cannot remove your own Administrator role.")
		}
	}

	c.Flash().Add("success", "User was updated successfully.")
	return c.Redirect(http.StatusSeeOther, "/admin/users")
}

// Destroy default implementation.
func (v AdminUsersResource) Destroy(c buffalo.Context) error {
	user, err := v.findUser(c)
	if err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	orgID := c.Session().Get("current_organization_id")

	// Instead of destroying the user, we remove them from the organization
	if err := tx.RawQuery("DELETE FROM user_roles WHERE user_id = ? AND organization_id = ?", user.ID, orgID).Exec(); err != nil {
		return err
	}

	c.Flash().Add("success", "User was removed from the organization.")
	return c.Redirect(http.StatusSeeOther, "/admin/users")
}

// PasswordReset allows admin to reset a user's password
func (v AdminUsersResource) PasswordReset(c buffalo.Context) error {
	currentUser := c.Value("current_user").(*models.User)
	isSuperAdmin := currentUser.HasRole("SUPER_ADMIN")

	if !isSuperAdmin {
		c.Flash().Add("danger", "Only Super Administrators can reset other users' passwords.")
		return c.Redirect(http.StatusSeeOther, "/admin/users")
	}

	user, err := v.findUser(c)
	if err != nil {
		return err
	}

	if err := user.GenerateResetToken(); err != nil {
		return err
	}

	tx := c.Value("tx").(*pop.Connection)
	if err := tx.Update(user); err != nil {
		return err
	}

	link := fmt.Sprintf("%spassword/edit/%s?token=%s", siteURL(), user.ID, user.ResetToken.String)
	w.Perform(worker.Job{
		Queue:   "default",
		Handler: "send_admin_password_reset",
		Args: worker.Args{
			"email": user.Email,
			"link":  link,
		},
	})

	c.Flash().Add("success", "Password reset link generated and sent to user.")
	return c.Redirect(http.StatusSeeOther, "/admin/users")
}
