# Pop ORM Best Practices & Gotchas

When writing or modifying models and database interactions in this project, AI agents must abide by the following Pop ORM specific rules.

## 1. Eager Loading Relationships
Pop ORM **does not automatically load** related models (like `has_many`, `belongs_to`, `many_to_many`).
If you query a user and expect their `Organizations` slice to be populated, you must explicitly tell Pop to load it using `.Eager()`.

**Example:**
```go
// WRONG: user.Organizations will be empty
err := tx.Find(&user, id)

// CORRECT: Loads all configured relationships
err := tx.Eager().Find(&user, id)

// CORRECT: Loads only specific relationships
err := tx.Eager("Organizations").Find(&user, id)
```

## 2. Using the Right Hooks
When generating fields programmatically (e.g., password hashing, generating unique slugs, or setting default UUIDs), do **not** use the `BeforeCreate` hook if those fields are marked with validation tags (e.g., `validate:"required"`).

Pop runs validations *before* `BeforeCreate`. If the field is empty, validation will fail.
Instead, use the `BeforeValidate` hook.

**Example:**
```go
// CORRECT: Generates the hash before Pop checks if PasswordHash is empty
func (u *User) BeforeValidate(tx *pop.Connection) error {
    if u.Password != "" {
        hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.PasswordHash = string(hash)
    }
    return nil
}
```

## 3. Custom Table Names
If your struct name doesn't neatly map to the pluralized table name, you can implement the `TableName` method:
```go
func (u User) TableName() string {
    return "app_users"
}
```
