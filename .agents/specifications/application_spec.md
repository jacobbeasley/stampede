# Application Specification: To-Do List Manager

## 1. Overview
This application is a task management platform that allows users to create accounts, log in, and securely manage their own personal to-do lists. The focus is on a clean, modern user interface backed by a robust and efficient server.

## 2. Technology Stack
- **Backend Framework:** [Buffalo](https://gobuffalo.io/) (Go web framework)
- **Frontend Reactive Framework:** [Svelte](https://svelte.dev/)
- **UI & Styling:** [DaisyUI](https://daisyui.com/) (Tailwind CSS component library)
- **Database:** PostgreSQL (Standard for Buffalo)

*Note: As requested, the architecture relies heavily on Buffalo (for routing, ORM, and server-side rendering) and DaisyUI (for the presentation layer), utilizing Svelte specifically for interactive client-side components.*

---

## 3. Data Model

The application uses a multi-tenant data model with **Organization**, **User**, **Role**, and **Todo** entities.

### `Organization`
Represents a tenant in the system.
| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `name` | String | Organization name |
| `created_at` | Timestamp | Record creation time |
| `updated_at` | Timestamp | Record last update time |

### `User`
Stores authentication, profile, and security tracking data.
| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `email` | String | User's email address (unique) |
| `first_name` | String | (Optional) User's first name |
| `last_name` | String | (Optional) User's last name |
| `phone_number` | String | (Optional) User's phone number |
| `password_hash` | String | Bcrypt hashed password |
| `reset_token` | String | (Optional) Token for password reset flow |
| `reset_token_expires_at` | Timestamp | (Optional) Expiration time for reset token |
| `pending_email` | String | (Optional) Temporary storage for unconfirmed new emails |
| `email_confirmation_token` | String | (Optional) Token for email update/signup confirmation |
| `email_confirmation_expires_at`| Timestamp | (Optional) Expiration time for email confirmation token |
| `failed_login_attempts` | Integer | Counter for failed logins (lockout threshold: 5) |
| `last_failed_login_at` | Timestamp | (Optional) Timestamp of the last failed login |
| `last_password_reset_request_at`| Timestamp | (Optional) Timestamp of the last password reset request (throttling: 5m) |
| `account_verified` | Boolean | Verification status (default: false) |
| `created_at` | Timestamp | Record creation time |
| `updated_at` | Timestamp | Record last update time |

### `Role`
Defines available roles in the system (`USER`, `ADMIN`, `SUPER_ADMIN`).
| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | String | Primary Key (e.g., 'ADMIN') |
| `role_name` | String | Display name |

### `UserRole` (Join Table)
Associates users with roles within a specific organization.
| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `user_id` | UUID | Foreign Key to `User` |
| `role_id` | String | Foreign Key to `Role` |
| `organization_id` | UUID | Foreign Key to `Organization` |

### `Todo`
Stores individual task items, linked to a user and an organization.
| Field | Type | Description |
| :--- | :--- | :--- |
| `id` | UUID | Primary Key |
| `user_id` | UUID | Foreign Key referencing `User.id` |
| `organization_id` | UUID | Foreign Key referencing `Organization.id` |
| `title` | String | Short title of the task |
| `description` | Text | (Optional) Detailed description of the task |
| `is_completed`| Boolean | Status of the task (default: false) |
| `position` | Integer | Sort order |
| `created_at` | Timestamp | Record creation time |
| `updated_at` | Timestamp | Record last update time |

---

## 4. Sitemap & Routing

### Public Routes
- `/` - **Landing Page**: Information about the app.
- `/register` - **Registration Page**: Form to create a new account (defaults to Default Organization).
- `/login` - **Login Page**: Form to authenticate an existing user.
- `/password/reset` - **Forgot Password Page**: Form to request a reset link.
- `/password/edit` - **Reset Password Page**: Form to set a new password using a token.
- `/api/ready` - **Liveness Probe**: Instant check if app has started.
- `/api/health` - **Readiness Probe**: Verifies database connectivity and app health.

### Authenticated Routes (Requires Session)
- `/auth/select_organization` - **Organization Selection**: For users with multiple organizations.
- `/todos` - **Dashboard (To-Do List)**: Scoped to the selected organization.
- `/logout` - **Logout Action**: Clears the session.

### Admin Routes (Requires Session & Admin Role)
- `/admin/users` - **User Management**: Manage users within the current organization.
- `/admin/users/new` - **Add User**: Invite or create a user for the current organization.
- `/admin/users/{id}/edit` - **Edit User**: Modify user roles within the organization.
- `/admin/users/{id}` (DELETE) - **Remove User**: Remove a user from the current organization.

### Super Admin Routes (Requires Session & Super Admin Role)
- `/admin/super/organizations` - **Organization Management**: List, create, and edit organizations.

*(Note: Most of the application will be standard full-stack Buffalo using HTML templates and DaisyUI. Svelte is used specifically on the `/todos` page to provide a highly interactive, reactive experience for managing tasks).*

---

## 5. Folder Layout & File Structure

The project will follow the standard Buffalo directory structure with additions to accommodate Svelte, Tailwind CSS, and DaisyUI in the frontend asset pipeline.

```text
buffalo-app/
├── actions/                 # Go Controllers and Routing logic
│   ├── app.go               # Main application routing and middleware setup
│   ├── auth.go              # Login, register, and session handling
│   ├── home.go              # Landing page handler
│   └── todos.go             # Todo CRUD operations (or JSON API for Svelte)
├── assets/                  # Frontend assets (bundled by Webpack/Vite)
│   ├── css/
│   │   └── application.css  # Tailwind and DaisyUI imports
│   ├── js/
│   │   ├── application.js   # Main frontend entry point (mounts Svelte)
│   │   └── components/      # Svelte components directory
│   │       ├── TodoApp.svelte
│   │       └── TodoItem.svelte
├── grifts/                  # Buffalo task scripts (e.g., database seeding)
├── locales/                 # i18n translation files
├── models/                  # Pop ORM models (Go structs for DB schema)
│   ├── user.go
│   └── todo.go
├── public/                  # Compiled/Static assets served directly
├── templates/               # Plush templates (Server-Side HTML)
│   ├── _flash.plush.html    # Flash messages template
│   ├── application.plush.html # Base HTML layout
│   ├── auth/                # Login/Register HTML pages
│   └── home/                # Landing HTML page
├── database.yml             # Pop database configuration
├── go.mod / go.sum          # Go module dependencies
├── package.json             # NPM dependencies (Svelte, Tailwind, DaisyUI)
└── tailwind.config.js       # Tailwind CSS & DaisyUI plugin configuration
```

---

## 6. Recommended Additional Sections

Beyond Data Model, Sitemap, and Folder Layout, you should consider including the following sections in a complete specification:

### A. Features & User Stories
Explicitly define what the user (and developer) can do to prevent scope creep.

**Developer / Infrastructure Stories:**
- *As a developer, I want to scaffold the initial Buffalo repository with Tailwind, DaisyUI, and Svelte integration so I have a solid foundation to build on.*
- *As a developer, I want to implement automated tests (unit and integration) so I can ensure the application remains stable as it grows.*
- *As a developer, I want comprehensive documentation (README, API specs, setup instructions) so new contributors can easily onboard.*

**End-User Stories:**
- *As a user, I want to sign up with my email so I can have my own private list.*
- *As a user, I want to mark a task as completed so I can track my progress.*
- *As a user, I want to edit a task's text in case I made a typo.*
- *As a user, I want to drag and drop tasks to reorder them in my list.*
- *As a user, I want to securely reset my password if I forget it.*
- *As a user, I want to manage my profile data such as name, email, and password while logged in.*

**Admin Stories:**
- *As an admin, I want a portal to manage (create, edit, delete, reset passwords) user accounts.*

### B. Internal API Endpoints (For Svelte Integration)
To support the Svelte application on the `/todos` page, Buffalo will expose the following JSON endpoints for task CRUD operations:
- `GET /api/todos` (Fetch user's to-dos, ordered by `position`)
- `POST /api/todos` (Create a to-do)
- `PUT /api/todos/{id}` (Update a to-do's text or completion status)
- `PUT /api/todos/reorder` (Batch update `position` fields after drag-and-drop)
- `DELETE /api/todos/{id}` (Delete a to-do)

### C. Security & Authentication
- **User Management**: We will use the `buffalo generate auth` command to scaffold the initial sign-up, login, logout, and session middleware.
- **Password Handling**: Passwords are hashed using bcrypt (included via the generator).
- **Session & CSRF**: Maintained via secure, HTTP-only cookies (included via the generator). Session rotation is applied on login state changes (login, password reset) to prevent session fixation. CSRF tokens must be extracted from the meta tag and sent as `authenticity_token` in POST/PUT/DELETE requests.
- **Password Reset Flow**: Throttled to once every 5 minutes per user. Generates a unique secure token, emails it asynchronously, and clears/expires the token upon password update.
- **Email Verification Flow**: Newly created accounts are unverified (`account_verified = false`) until they verify via the email token link. Email updates on the profile page trigger a verification workflow sent to the old email before updating to the new one.
- **Rate Limiting & Lockout**: Implemented on the login flow (5 failed attempts lock the account for 5 minutes).
- **Authorization Middleware**: Rules to ensure users can only access rows in the `todos` table that match their `user_id` and the active tenant `organization_id`.
- **Admin Middleware**: Rules to restrict access to `/admin` routes solely to users with the `ADMIN` or `SUPER_ADMIN` role inside their organization.
- **Super Admin Middleware**: Rules to restrict organization creation/management solely to global `SUPER_ADMIN` users.

### D. UI/UX Design System
- State the intention to use DaisyUI themes (e.g., specifying a default dark/light theme).
- Mention any required micro-animations or specific responsive behaviors.

### E. Testing Strategy
- Define the expectations for testing (e.g., unit tests for the Go `models/`, integration tests for `actions/`, and potentially component tests for Svelte).

### F. Environment & Deployment
- **Environment Variables**: Define necessary config values like `DATABASE_URL`, `SESSION_SECRET`, and `SMTP_HOST` (for password reset emails).
- **Deployment**: Outline the hosting strategy (e.g., building a Docker container via Buffalo's standard Dockerfile and deploying to a service like Render or AWS).

### G. Database Optimization
- **Indexing**: Ensure indexes are created in PostgreSQL on frequently queried fields, particularly the `user_id` on the `todos` table (since queries will always filter by user) and `email` on the `users` table.

### H. Documentation
- **Developer Documentation**: A comprehensive `README.md` covering local setup, environment variables, database seeding, and testing instructions.
- **API Documentation**: Clear documentation (e.g., Swagger/OpenAPI or a simple Markdown file) for the internal JSON endpoints used by the Svelte frontend.
