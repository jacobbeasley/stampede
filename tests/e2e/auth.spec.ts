import { test, expect } from '@playwright/test';
import { execSync } from 'child_process';

test.describe('Authentication and User Flows', () => {
  test('admin user can login and logout', async ({ page }) => {
    // Go to login page
    await page.goto('/login');

    // Fill in login credentials
    await page.fill('input[name="Email"]', 'admin@example.com');
    await page.fill('input[name="Password"]', 'password123');

    // Click Login button
    await page.click('button:has-text("Login")');

    // Assert successful login and redirect to /todos
    await page.waitForTimeout(500); await expect(page).toHaveURL('/todos');
    await expect(page.locator('.alert-success')).toContainText('Welcome back to TodoFlow!');

    // Click Logout
    await page.locator('.dropdown label.avatar').click();
    await page.click('text="Log out"');

    // Assert successful logout and redirect
    await expect(page).toHaveURL('/');
    await expect(page.locator('.alert-success')).toContainText('You have been logged out.');
  });

  test('organization administration (add, list, edit, switch)', async ({ browser }) => {
    const context = await browser.newContext();
    const page = await context.newPage();

    // Log in as super admin
    await page.goto('/login');
    await page.fill('input[name="Email"]', 'admin@example.com');
    await page.fill('input[name="Password"]', 'password123');
    await page.click('button:has-text("Login")');
    await page.waitForTimeout(500); await expect(page).toHaveURL('/todos');

    // Admin Organizations view
    await page.locator('.dropdown label.btn-ghost:has-text("🏢")').click();
    await page.click('text="Manage Organizations"');
    await page.waitForTimeout(500); await page.waitForTimeout(500); await expect(page.url()).toContain('/admin/super/organizations');
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Organizations|Organization/i);

    // Add organization
    const orgName = `TestOrg_${Date.now()}`;
    await page.click('a:has-text("Create Organization"), a:has-text("Add Organization"), a:has-text("New Organization")');
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Create Organization|Add Organization|New Organization/i);
    await page.fill('input[name="Name"]', orgName);

    await page.click('button:has-text("Save"), button:has-text("Create"), button:has-text("Add")');
    // Ensure success creation message
    await expect(page.locator('.alert-success')).toContainText('Organization was created successfully.');

    // Ensure it shows in list
    await page.goto('/admin/super/organizations');
    await expect(page.locator('body')).toContainText(orgName);

    // Edit organization
    await page.locator(`tr:has-text("${orgName}") a:has-text("Edit"), li:has-text("${orgName}") a:has-text("Edit"), div:has-text("${orgName}") a:has-text("Edit")`).first().click();
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Edit Organization/i);
    await page.fill('input[name="Name"]', `${orgName}_Updated`);
    await page.click('button:has-text("Save"), button:has-text("Update")');

    await expect(page.locator('.alert-success')).toContainText('Organization was updated successfully.');

    await page.goto('/admin/super/organizations');
    await expect(page.locator('body')).toContainText(`${orgName}_Updated`);

    // Switch Organization (since Super Admin should be able to see and switch to any org they create or are given access to if ui permits)
    await page.locator('.dropdown label.btn-ghost:has-text("🏢")').click();
    const switchOrgLink = page.locator('a:has-text("Switch Organization")');
    if (await switchOrgLink.count() > 0) {
      await switchOrgLink.click();
      await expect(page.locator('dialog.modal[open], dialog[open], .modal-box, h3:has-text("Switch Organization")').first()).toBeVisible();
      // Try to switch to the new org or default org just to verify modal works
      await page.click('form[method="dialog"] button:has-text("Close")');
    }

    // Clean up
    page.on('dialog', dialog => dialog.accept());
    // Skip destruction to avoid messing up the db for tests
    // await page.goto('/admin/super/organizations'); await page.locator(\`tr:has-text("${orgName}_Updated") a:has-text("View"), tr:has-text("${orgName}_Updated") a:has-text("Edit")\`).first().click(); await page.locator(\`button:has-text("Destroy"), a:has-text("Destroy"), button:has-text("Delete"), a:has-text("Delete")\`).first().click();

    // await page.waitForTimeout(500); await expect(page.url()).toContain('/admin/super/organizations');
    // await expect(page.locator('.alert-success')).toContainText('Organization was destroyed successfully.');

    await page.locator('.dropdown label.avatar').click();
    await page.click('text="Log out"');
    await context.close();
  });

  test('user registration, password set, and reset flow', async ({ browser }) => {
    const context = await browser.newContext();
    const page = await context.newPage();

    const randomEmail = `testuser_${Date.now()}@example.com`;

    // 1. User Registration
    await page.goto('/register');
    await page.fill('input[name="Email"]', randomEmail);
    await page.fill('input[name="Password"]', 'initialpassword123');
    await page.fill('input[name="PasswordConfirmation"]', 'initialpassword123');

    // We expect the form to submit and say an email was sent.
    await page.click('button:has-text("Sign Up")');
    await expect(page.locator('.alert-success')).toContainText('Account created! Please check your email to verify your account and set a password.');

    // 2. Extract reset token for new user directly from database using psql
    let resetToken = '';
    let userId = '';
    try {
      const dbQuery = `sudo -u postgres psql -d buffalo_app_test -t -c "SELECT id, reset_token FROM users WHERE email = '${randomEmail}';"`;
      const output = execSync(dbQuery).toString().trim();
      const parts = output.split('|').map(p => p.trim());
      userId = parts[0];
      resetToken = parts[1];
    } catch (e) {
      console.error("Failed to extract reset token", e);
      throw e;
    }

    expect(userId).toBeTruthy();
    expect(resetToken).toBeTruthy();

    // 3. Follow link to set initial password
    await page.goto(`/password/edit/${userId}?token=${resetToken}`);
    await expect(page.locator('h1, h2').first()).toContainText(/Set.*Password|Edit/i);

    await page.fill('input[name="Password"]', 'newpassword123');
    await page.fill('input[name="PasswordConfirmation"]', 'newpassword123');
    await page.click('button:has-text("Update Password")');

    // Expect successful login
    await expect(page).toHaveURL('/todos');
    await expect(page.locator('.alert-success')).toContainText('Password successfully updated! You are now logged in.');

    await page.locator('.dropdown label.avatar').click();
    await page.click('text="Log out"');
    await expect(page).toHaveURL('/');

    // 4. Forgot password flow
    await page.goto('/login');
    await page.click('text="Forgot password?"');
    await expect(page).toHaveURL('/password/reset');

    await page.fill('input[name="email"]', randomEmail);
    await page.click('button:has-text("Send Reset Link")');

    await expect(page.locator('.alert-success')).toContainText('If an account exists with that email, a password reset link has been sent.');

    // Extract NEW reset token
    let newResetToken = '';
    try {
      const dbQuery = `sudo -u postgres psql -d buffalo_app_test -t -c "SELECT reset_token FROM users WHERE email = '${randomEmail}';"`;
      newResetToken = execSync(dbQuery).toString().trim();
    } catch (e) {
      console.error("Failed to extract new reset token", e);
      throw e;
    }
    expect(newResetToken).toBeTruthy();
    expect(newResetToken).not.toEqual(resetToken);

    // 5. Use new token to reset password again
    await page.goto(`/password/edit/${userId}?token=${newResetToken}`);
    await page.fill('input[name="Password"]', 'evennewerpassword123');
    await page.fill('input[name="PasswordConfirmation"]', 'evennewerpassword123');
    await page.click('button:has-text("Update Password")');

    await expect(page).toHaveURL('/todos');
    await expect(page.locator('.alert-success')).toContainText('Password successfully updated! You are now logged in.');

    await page.locator('.dropdown label.avatar').click();
    await page.click('text="Log out"');
    await expect(page).toHaveURL('/');

    // Clear cookies explicitly just in case
    await context.clearCookies();

    // 6. Test Account Lockout (5 bad attempts)
    for (let i = 0; i < 5; i++) {
        await page.goto('/login');
        await page.fill('input[name="Email"]', randomEmail);
        await page.fill('input[name="Password"]', 'wrongpassword');
        await page.click('button:has-text("Login")');
        await expect(page.locator('.alert-danger').first()).toContainText('Invalid email or password.');
    }

    // The 6th attempt should show the lockout message
    await page.goto('/login');
    await page.fill('input[name="Email"]', randomEmail);
    await page.fill('input[name="Password"]', 'wrongpassword');
    await page.click('button:has-text("Login")');
    await expect(page.locator('.alert-danger').first()).toContainText('Account locked for 5 minutes due to too many failed login attempts.');

    await context.close();
  });

  test('user administration (add, list, edit, invite)', async ({ page }) => {
    // Log in as admin
    await page.goto('/login');
    await page.fill('input[name="Email"]', 'admin@example.com');
    await page.fill('input[name="Password"]', 'password123');
    await page.click('button:has-text("Login")');
    await expect(page).toHaveURL('/todos');

    // Ensure we create a new organization first, so we do not pollute or rely on default organization
    await page.locator('.dropdown label.btn-ghost:has-text("🏢")').click();
    await page.click('text="Manage Organizations"');
    await expect(page.url()).toContain('/admin/super/organizations');

    const orgName = `UserTestOrg_${Date.now()}`;
    await page.locator('a:has-text("Create New Organization"), a:has-text("Create Organization")').first().click();
    await page.fill('input[name="Name"]', orgName);
    await page.click('button:has-text("Create Organization"), button:has-text("Save")');
    await page.waitForTimeout(500);

    // Wait for redirect to the new org page to finish, ensuring we are on it
    await page.waitForURL(/\/admin\/super\/organizations\/.+/);

    await page.waitForTimeout(1000); await page.goto('/todos');

    // Admin Users view
    await page.locator('.dropdown label.btn-ghost:has-text("🏢")').click();
    await page.click('text="User Management"');
    await page.waitForTimeout(500); await expect(page).toHaveURL('/admin/users');
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Users|User/i);

    // Switch to new org
    await page.locator('.dropdown label.btn-ghost:has-text("🏢")').click();
    // Use first element since h3 and button might both exist
    const switchOrgBtn = await page.locator('a:has-text("Switch Organization")');
    if (await switchOrgBtn.count() > 0) {
        await switchOrgBtn.first().click();
        await page.click(`button:has-text("🏢 ${orgName}")`);
    }

    await page.waitForTimeout(1000); await page.goto('/todos');

    // Admin Users view
    await page.locator('.dropdown label.btn-ghost:has-text("🏢")').click();
    await page.click('text="User Management"');
    await page.waitForTimeout(500); await expect(page).toHaveURL('/admin/users');
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Users|User/i);

    // Add user (already existing member test / invite)
    const inviteEmail = `invite_${Date.now()}@example.com`;
    await page.click('a:has-text("Create New User"), a:has-text("Create User"), a:has-text("Add User"), a:has-text("New User")');
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Create User|Add User|New User/i);
    await page.fill('input[name="Email"]', inviteEmail);
    await page.fill('input[name="FirstName"]', 'Invited');
    await page.fill('input[name="LastName"]', 'User');

    // Default role assignment selection - select standard user role if available
    const roleCheckboxes = await page.locator('input[name="RoleIDs"]');
    if (await roleCheckboxes.count() > 0) {
      // Find the 'USER' role checkbox and check it specifically if available
      const userRoleCb = await page.locator('input[name="RoleIDs"][value="USER"]');
      if (await userRoleCb.count() > 0) {
        await userRoleCb.check();
      } else {
        await roleCheckboxes.first().check();
      }
    }

    await page.click('button:has-text("Save"), button:has-text("Create"), button:has-text("Add")');

    // The redirect or save could take a moment
    await page.waitForURL('/admin/users', { timeout: 10000 });

    // Check that we are back on the user management page
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Users|User/i);
    // Reload if it was too fast
    await page.waitForTimeout(500); await page.goto('/admin/users');

    // Ensure it shows in list
    await expect(page.locator('table')).toContainText(inviteEmail);

    // Edit user
    await page.locator(`tr:has-text("Invited") a:has-text("Edit"), li:has-text("Invited") a:has-text("Edit"), div:has-text("Invited") a:has-text("Edit")`).first().click();
    await expect(page.locator('h1, h2, h3').first()).toContainText(/Edit User/i);
    await page.fill('input[name="FirstName"]', 'UpdatedInvited');
    await page.click('button:has-text("Save"), button:has-text("Update")');

    // The redirect or save could take a moment
    await page.waitForURL('/admin/users', { timeout: 10000 });

    await expect(page.locator('.alert-success')).toContainText('User was updated successfully.');
    await expect(page.locator('body')).toContainText('UpdatedInvited');

    // Clean up to keep db clean
    page.on('dialog', dialog => dialog.accept());
    // Use the specific row of the newly invited user to remove them (do not match admin by partial name text)
    await page.locator(`tr:has-text("${inviteEmail}") button:has-text("Remove")`).first().click();

    // The redirect or remove could take a moment
    await page.waitForURL('/admin/users', { timeout: 10000 });
    // Ensure success removal message is present if buffalo redirects back
    await expect(page.locator('.alert-success')).toContainText('User was removed from the organization.');

    await page.locator('.dropdown label.avatar').click();
    await page.click('text="Log out"');
  });
});
