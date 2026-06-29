package actions

import (
	"fmt"
	"net/http"

	"buffalo-app/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
)

// TodosIndex renders the main page that hosts the Svelte TodoApp
func TodosIndex(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("todos/index.plush.html"))
}

// TodosList gets the current user's todos
func TodosList(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(*models.User)
	orgID := c.Session().Get("current_organization_id")

	todos := &models.Todoes{}
	// Get todos for current user, ordered by position
	if err := tx.Where("user_id = ? AND organization_id = ?", user.ID, orgID).Order("position ASC").All(todos); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(todos))
}

// TodosCreate creates a new todo for the current user
func TodosCreate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(*models.User)
	orgIDRaw := c.Session().Get("current_organization_id")
	orgID, err := uuid.FromString(fmt.Sprintf("%v", orgIDRaw))
	if err != nil {
		return err
	}

	todo := &models.Todo{}
	if err := c.Bind(todo); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	todo.UserID = user.ID
	todo.OrganizationID = orgID

	// If position is not explicitly set, we might want to put it at the end
	// Let's find the max position for the user
	if todo.Position == 0 {
		var maxPos int
		err := tx.RawQuery("SELECT COALESCE(MAX(position), 0) FROM todoes WHERE user_id = ? AND organization_id = ?", user.ID, orgID).First(&maxPos)
		if err == nil {
			todo.Position = maxPos + 1
		}
	}

	verrs, err := tx.ValidateAndCreate(todo)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusCreated, r.JSON(todo))
}

// TodosUpdate updates an existing todo
func TodosUpdate(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(*models.User)
	orgID := c.Session().Get("current_organization_id")

	todo := &models.Todo{}
	if err := tx.Where("id = ? AND user_id = ? AND organization_id = ?", c.Param("todo_id"), user.ID, orgID).First(todo); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("todo not found"))
	}

	if err := c.Bind(todo); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	// Make sure user ID cannot be changed
	todo.UserID = user.ID

	verrs, err := tx.ValidateAndUpdate(todo)
	if err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	if verrs.HasAny() {
		return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
	}

	return c.Render(http.StatusOK, r.JSON(todo))
}

// TodosDelete deletes a todo
func TodosDelete(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(*models.User)
	orgID := c.Session().Get("current_organization_id")

	todo := &models.Todo{}
	if err := tx.Where("id = ? AND user_id = ? AND organization_id = ?", c.Param("todo_id"), user.ID, orgID).First(todo); err != nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("todo not found"))
	}

	if err := tx.Destroy(todo); err != nil {
		return c.Error(http.StatusInternalServerError, err)
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "deleted successfully"}))
}

// TodoReorderRequest is the payload for batch reordering
type TodoReorderRequest struct {
	Order []string `json:"order"`
}

// TodosReorder updates positions for multiple todos
func TodosReorder(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	user := c.Value("current_user").(*models.User)
	orgID := c.Session().Get("current_organization_id")

	req := &TodoReorderRequest{}
	if err := c.Bind(req); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}

	for i, id := range req.Order {
		todo := &models.Todo{}
		if err := tx.Where("id = ? AND user_id = ? AND organization_id = ?", id, user.ID, orgID).First(todo); err != nil {
			// Skip invalid IDs
			continue
		}
		todo.Position = i + 1
		// We only want to update position
		if err := tx.Update(todo, "position"); err != nil {
			return c.Error(http.StatusInternalServerError, err)
		}
	}

	return c.Render(http.StatusOK, r.JSON(map[string]string{"message": "reordered successfully"}))
}
