---
name: implement-sso
description: Guides the fullstack-builder on how to implement Single Sign-On (SSO) alongside existing local authentication in the Buffalo framework, avoiding destructive CLI overwrites and safely linking user data models.
---

# Implement SSO Skill

This skill outlines how to enable Single Sign-On (SSO) in this Buffalo project WITHOUT destroying the existing local email/password login system.

## 1. Context and Planning
Buffalo uses the `goth` package (`github.com/markbates/goth`) to manage OAuth safely.
- **CRITICAL**: The application already has local authentication (`actions/auth.go`) and a `models.User` schema.
- **DO NOT** run the standard `buffalo generate goth {provider}` command, as it will prompt to overwrite and destroy the existing `actions/auth.go` file and local auth logic.

## 2. Safe Generation & Integration
Instead of full boilerplate generation, use the safe sub-generator or manual integration:
- Run `buffalo generate goth-auth {provider}` (e.g., `buffalo generate goth-auth google`). This only generates necessary routing and provider allocation.
- **Manual Merge**: If the generator creates a separate file or prompts for overwrite, manually merge the Goth provider initialization into your `actions/app.go` or `actions/auth.go`:
  ```go
  goth.UseProviders(
      google.New(envy.Get("GOOGLE_KEY", ""), envy.Get("GOOGLE_SECRET", ""), fmt.Sprintf("%s/auth/google/callback", envy.Get("APP_URL", "http://localhost:3000"))),
  )
  ```
- Map the routes manually in `actions/app.go`:
  - `app.GET("/auth/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))`
  - `app.GET("/auth/{provider}/callback", AuthCallback)`

## 3. Data Model & Account Linking
The existing `models.User` requires `Email` and `PasswordHash`. OAuth users will not have a password natively.
In your `AuthCallback` handler:
1. Extract the `goth.User` from the request context: `user, err := goth_buffalo.UserFromContext(c)`.
2. **Account Linking**: Look up the user in the database by their `goth.User.Email`.
   - **If they exist**: Log them in immediately. You have successfully linked their OAuth login to their existing local account.
   - **If they DO NOT exist**: Instantiate a new `models.User`.
     - Map `goth.User.Email` to `User.Email`.
     - Map `goth.User.FirstName` and `LastName`.
     - **Bypass Constraints**: Generate a random secure string (e.g., using `crypto/rand` or `uuid.NewV4()`) and assign it to `User.Password` and `User.PasswordConfirmation` so the model passes its `BeforeValidate` hook requiring a password hash.
     - Save the user: `err := tx.Create(&newUser)`.
3. Create a new Gorilla session for the user.
4. Redirect to the dashboard.

## 4. Frontend Integration
- On the Svelte 5 login page, add the SSO options *alongside* the existing local login form. Provide standard HTML links to the Buffalo auth routes, not AJAX fetches.
  - `<a href="/auth/google" class="btn btn-outline">Sign in with Google</a>`

## 5. Security Rules
- **Session Fixation**: As defined in `AGENTS.md`, you MUST rotate the Gorilla session by expiring the old one (`c.Session().Session.Options.MaxAge = -1` and `Save()`), then clearing its properties (`ID = ""`, `IsNew = true`, and reinitializing `Values`) upon successful login.
- Never log or store the provider secret keys in code. Ensure they are loaded via `envy.Get()`.
