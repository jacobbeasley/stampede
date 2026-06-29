---
name: manage-database-schema
description: Use this skill when you need to make changes to the database schema, including creating new tables, adding columns, or updating relationships. It ensures safe migration practices using Buffalo/Soda CLIs and proper Pop ORM updates.
---

# Manage Database Schema

This skill outlines the necessary steps for safely modifying the database schema and Pop ORM models.

## 1. Context and Planning
Before changing the database schema:
- Identify the required fields, data types, and relationships (e.g., `has_many`, `belongs_to`).
- Determine whether to generate a new model or modify an existing table via a Fizz migration.
- Review existing models in `models/` to ensure naming consistency and avoid duplication.
- When creating new tables or modifying existing ones, always consider indexing. Use the `index-database` skill to ensure proper indexes are added for fields used in `WHERE`, `JOIN`, and `ORDER BY` clauses.

## 2. Generating Migrations
Never manually edit the database schema. Always use the provided tools:
- **Immutability Rule**: DO NOT modify existing migrations that have already been committed to the repository. Migrations are immutable. If you need to change the schema, you MUST generate a new migration.
- **New Models**: Use `buffalo pop generate model <ModelName> [fields...]` or `soda generate model <ModelName> [fields...]`. This generates both the Fizz migration and the Go model struct.
- **Modifying Existing Tables**: Generate a blank Fizz migration using `buffalo pop generate fizz <migration_name>` or `soda generate fizz <migration_name>`, and write the appropriate `add_column`, `drop_column`, or `change_column` commands.

## 3. Applying Migrations
Apply your changes to both the development and test databases:
- Ensure PostgreSQL is running.
- Run migrations: `buffalo pop migrate` or `soda migrate`.
- *Fallback*: If the `soda` CLI is missing, install it: `go install github.com/gobuffalo/pop/v6/soda@latest`.

## 4. Updating Pop ORM Models
After the schema has migrated:
- **Eager Loading**: Remember that Pop ORM does *not* auto-load relationships. If your new schema introduces relationships, update queries to use `.Eager()` where necessary (e.g., `tx.Eager().Find(&user, id)`).
- **Hooks**: Use the `BeforeValidate` hook to generate required fields programmatically (e.g., password hashing, default UUIDs) before Pop runs validation checks. Do *not* use `BeforeCreate` for fields that are marked as required in the struct tags.
- **Validation**: Update the `Validate`, `ValidateCreate`, and `ValidateUpdate` methods in your model to enforce business logic on the new fields.

## 5. Testing
Ensure database changes are covered by tests:
- **Database Preparation**: Tests require the test database to be fully migrated. Run `buffalo pop create -a` and `buffalo pop migrate` (or soda equivalents) before executing tests.
- **Model Tests**: Create or update the corresponding `_test.go` file in the `models/` directory to verify validation rules, hooks, and relationships. Run `buffalo test ./models/...` to confirm.
