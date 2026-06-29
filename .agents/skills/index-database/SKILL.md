---
name: index-database
description: Use this skill to evaluate a table and its queries to ensure proper database indexing for efficiency and performance.
---

# Index Database

This skill outlines the steps to examine a table and its queries to determine where indexes are needed.

## 1. Examine the Table and Queries
- Search the codebase (e.g., using `grep`) for database queries (`tx.Where`, `tx.RawQuery`, `tx.Eager`, etc.) targeting the table.
- Identify fields used in `WHERE` and `JOIN` clauses.
- Identify fields used for sorting (`ORDER BY` or `Order`).
- Look at foreign key relationships, as PostgreSQL does not automatically index foreign keys.

## 2. Identify Needed Indexes
- **Single-column indexes**: Needed for columns frequently used in `WHERE` clauses (e.g., `email`, `reset_token`) or foreign keys used in joins.
- **Composite indexes**: Needed when queries frequently filter on multiple columns together (e.g., `user_id` and `organization_id`). Order the columns based on selectivity and query patterns.
- **Unique indexes**: Needed to enforce uniqueness constraints (if not already handled by a primary key or existing unique index).

## 3. Generate Migrations
- NEVER modify existing migrations. Always generate a new migration.
- Generate a new migration for the indexes: `buffalo pop generate fizz add_indexes_to_<table_name>` or `soda generate fizz add_indexes_to_<table_name>`.
- Write the `add_index` commands in the `.up.fizz` file. Examples:
  - `add_index("users", "reset_token", {})`
  - `add_index("user_roles", ["user_id", "organization_id"], {})`
- Write the corresponding `drop_index` commands in the `.down.fizz` file. Examples:
  - `drop_index("users", "users_reset_token_idx")`

## 4. Apply and Verify
- Apply the migrations to your local development and test databases using `soda migrate` or `buffalo pop migrate`.
- Verify that the application functions correctly with the new indexes.
