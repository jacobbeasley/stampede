---
name: manage-mailers
description: Guides the creation and delivery of emails using Buffalo's built-in mailers and Plush templates, specifically integrated with the worker system.
---

# Manage Mailers Skill

This skill outlines how to generate, format, and dispatch emails in the Buffalo framework safely and efficiently.

## 1. Generating a Mailer
Buffalo provides generators for mailers. Use them to maintain convention.
- Run `buffalo generate mailer <name>` (e.g., `buffalo generate mailer welcome_email`).
- This will generate:
  - A Go file in the `mailers/` directory containing the logic to build the email.
  - An HTML Plush template in `templates/mail/`.
  - A plain-text Plush template in `templates/mail/`.

## 2. Formatting Email Templates (Plush)
- **HTML Templates**: Use the `.plush.html` files generated in `templates/mail/`.
  - Ensure the template uses inline CSS or references the centralized `templates/mail/layout.plush.html` where styles are defined.
  - Example: `<p>Welcome, <%= user.Name %>!</p>`
- **Plain Text Templates**: Use the `.plush.txt` files for clients that don't support HTML. Provide the same information without markup.

## 3. Building the Message
In your generated `mailers/` Go file, populate the `mail.Message`:
```go
func SendWelcomeEmail(user models.User) error {
    m := mail.NewMessage()
    m.From = "no-reply@example.com"
    m.To = []string{user.Email}
    m.Subject = "Welcome to Our Platform!"

    // Pass data to the Plush templates
    data := map[string]interface{}{
        "user": user,
    }

    // Add bodies (HTML and Plain text)
    err := m.AddBodies(data, "mail/welcome_email.plush.html", "mail/welcome_email.plush.txt")
    if err != nil {
        return err
    }

    return smtp.Send(m)
}
```

## 4. Dispatching via Background Worker (Crucial)
Sending emails synchronously during an HTTP request blocks the response. **You must dispatch emails via a background worker.**
- Do NOT call `SendWelcomeEmail()` directly inside a Buffalo Action.
- Follow the `.agents/skills/manage-background-jobs/SKILL.md` to create a worker handler.
- The worker handler should fetch necessary data (e.g., the User record) and then call the Mailer function.

## 5. Configuration
- Ensure SMTP credentials (`SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASSWORD`) are correctly read via `envy.Get` in `mailers/mailers.go`.
- In development/testing environments, Buffalo mailers fallback to a custom logger for console output if unconfigured.
