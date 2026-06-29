package actions

import (
	"net/http"

	"buffalo-app/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

type OrganizationsResource struct {
	buffalo.Resource
}

// List gets all Organizations. This is for super admins.
func (v OrganizationsResource) List(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	organizations := &models.Organizations{}
	q := tx.PaginateFromParams(c.Params())
	if err := q.All(organizations); err != nil {
		return err
	}
	c.Set("pagination", q.Paginator)
	c.Set("organizations", organizations)
	return c.Render(http.StatusOK, r.HTML("admin/organizations/index.plush.html"))
}

// Show gets a specific Organization.
func (v OrganizationsResource) Show(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	organization := &models.Organization{}
	if err := tx.Find(organization, c.Param("organization_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	c.Set("organization", organization)
	return c.Render(http.StatusOK, r.HTML("admin/organizations/show.plush.html"))
}

// New renders the form for creating a new Organization.
func (v OrganizationsResource) New(c buffalo.Context) error {
	c.Set("organization", &models.Organization{})
	return c.Render(http.StatusOK, r.HTML("admin/organizations/new.plush.html"))
}

// Create adds an Organization to the DB.
func (v OrganizationsResource) Create(c buffalo.Context) error {
	organization := &models.Organization{}
	if err := c.Bind(organization); err != nil {
		return err
	}
	tx := c.Value("tx").(*pop.Connection)
	verrs, err := tx.ValidateAndCreate(organization)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("organization", organization)
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/organizations/new.plush.html"))
	}
	c.Flash().Add("success", "Organization was created successfully.")
	return c.Redirect(http.StatusSeeOther, "/admin/super/organizations/%v", organization.ID)
}

// Edit renders the form for editing an Organization.
func (v OrganizationsResource) Edit(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	organization := &models.Organization{}
	if err := tx.Find(organization, c.Param("organization_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	c.Set("organization", organization)
	return c.Render(http.StatusOK, r.HTML("admin/organizations/edit.plush.html"))
}

// Update changes an Organization in the DB.
func (v OrganizationsResource) Update(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	organization := &models.Organization{}
	if err := tx.Find(organization, c.Param("organization_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if err := c.Bind(organization); err != nil {
		return err
	}
	verrs, err := tx.ValidateAndUpdate(organization)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("organization", organization)
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.HTML("admin/organizations/edit.plush.html"))
	}
	c.Flash().Add("success", "Organization was updated successfully.")
	return c.Redirect(http.StatusSeeOther, "/admin/super/organizations/%v", organization.ID)
}

// Destroy deletes an Organization from the DB.
func (v OrganizationsResource) Destroy(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	organization := &models.Organization{}
	if err := tx.Find(organization, c.Param("organization_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}
	if err := tx.Destroy(organization); err != nil {
		return err
	}
	c.Flash().Add("success", "Organization was destroyed successfully.")
	return c.Redirect(http.StatusSeeOther, "/admin/super/organizations")
}
