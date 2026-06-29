# Stampede: A Modern Buffalo Hackathon Starter

A multi-tenant, full-stack hackathon starter application built with **Buffalo** (Go), **Svelte 5**, **Tailwind CSS v4**, and **DaisyUI**. This starter is designed for an agentic workflow in which you use its skills to design landing pages and mockups for your site, then it can help build systems architectural specifications, API and database schema design, and then generate the initial implementation for your features (code, migrations, etc.). It even includes some documentation about how to deploy it on Google Cloud, as well as other platforms.

![Stampede](stampede.png)

### Demo

![demo](demo.gif)

### Key Features

- **Multi-tenancy**: Users can belong to multiple organizations and switch between them seamlessly. 
- **Role-Based Access Control (RBAC)**:
  - **Super Administrators**: Manage organizations and all users across the platform.
  - **Administrators**: Manage users within their specific organization.
  - **Users**: Manage their personal to-do lists within an organization.
- **Sample Feature**:Ships with a to-do example feature in which todo lists are strictly scoped to both a user and an organization.
- **Easy Onboarding**: New users can sign up and are automatically joined to the Default Organization.
- **Efficient Runtime**: Compiles to 25mb go executable (30mb docker image) and runs with no Javascript frontend dependencies (compiles to native JS with Svelte).

### Quick Start

Get the application up and running locally in three quick steps:

1. **Install Dependencies & Set Up Environment**:
   ```bash
   # Install backend & frontend dependencies
   go mod download
   npm install
   go install github.com/gobuffalo/cli/cmd/buffalo@latest
   go install github.com/gobuffalo/buffalo-pop/v3@latest

   # Copy the sample environment file
   cp .env.sample .env
   ```

2. **Prepare the Database**:
   Ensure **PostgreSQL** is running locally, then initialize and migrate the database:
   ```bash
   buffalo pop create -a
   buffalo pop migrate
   ```

3. **Start the Development Servers**:
   Launch the Go backend and Svelte/Vite asset compilation in separate terminals:
   ```bash
   # Terminal 1: Backend Dev Server (runs on http://localhost:3000)
   buffalo dev

   # Terminal 2: Frontend Asset Pipeline (runs on http://localhost:3001)
   npm run dev
   ```

Now open **[http://localhost:3000](http://localhost:3000)** in your browser!

> [!TIP]
> **Default Admin Account:** Log in with the pre-seeded Super Administrator credentials:
> - **Username:** `admin@example.com`
> - **Password:** `password123`

---

## Table of Contents

- [Tech Stack](#tech-stack)
- [Repository Structure](#repository-structure)
- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Local Development (Without Docker)](#local-development-without-docker)
  - [Linux & macOS](#linux--macos)
  - [Windows](#windows)
- [Running with Docker](#running-with-docker)
- [Database Management](#database-management)
- [Building for Production](#building-for-production)
- [Stopping the Application](#stopping-the-application)
- [Running Tests](#running-tests)

---

## Tech Stack

| Layer | Technology | Version |
|---|---|---|
| Backend framework | [Buffalo](https://gobuffalo.io/) | v1.1.4 |
| Language | Go | 1.26+ |
| ORM | [Pop](https://gobuffalo.io/documentation/database/pop/) (gobuffalo/pop/v6) | v6.3 |
| Database | PostgreSQL | 14+ |
| Template engine | [Plush](https://github.com/gobuffalo/plush) | v5 |
| Frontend framework | [Svelte](https://svelte.dev/) | v5 (Runes syntax) |
| Frontend build tool | [Vite](https://vitejs.dev/) | v8 |
| CSS framework | [Tailwind CSS](https://tailwindcss.com/) | v4 |
| UI component library | [DaisyUI](https://daisyui.com/) | v5 |
| Session security | [unrolled/secure](https://github.com/unrolled/secure) | v1 |

---

## Repository Structure

```text
buffalo-app/
├── actions/                    # Go controllers (thin handlers — push logic to models/)
│   ├── app.go                  # Route definitions and middleware setup
│   ├── home.go                 # Landing page handler
│   ├── render.go               # Buffalo render engine configuration
│   └── vite.go                 # Vite manifest helper (resolves hashed asset filenames)
├── assets/                     # Frontend source files (compiled by Vite, NOT served directly)
│   ├── css/
│   │   └── application.css     # Tailwind CSS v4 entry point + DaisyUI plugin
│   └── js/
│       ├── main.js             # Vite/Svelte entry point (mounts Svelte into #app)
│       └── components/         # Svelte components
│           └── HelloWorld.svelte
├── cmd/app/
│   └── main.go                 # Application entrypoint
├── config/
│   └── buffalo-app.toml        # Buffalo application metadata
├── grifts/                     # Buffalo task scripts (e.g. database seeding)
├── locales/                    # i18n translation files
├── migrations/                 # Database migration files (Fizz format)
├── models/                     # Pop ORM models (Go structs matching DB schema)
│   └── models.go               # DB connection setup
├── public/                     # Statically served files
│   └── assets/                 # Compiled frontend output (generated by `npm run build`)
├── templates/                  # Plush server-rendered HTML templates
│   ├── application.plush.html  # Base HTML layout (navbar, footer, asset tags)
│   ├── _flash.plush.html       # Flash message partial
│   └── home/
│       └── index.plush.html    # Landing page
├── .agents/                    # AI agent configuration (do not delete)
│   ├── specifications/         # Application spec and story documentation
│   │   ├── application_spec.md
│   │   └── stories/            # Backlog and per-story folders
│   ├── skills/                 # Reusable AI agent skills
│   └── knowledge/              # Knowledge items (KIs) about the codebase
├── .buffalo.dev.yml            # Buffalo hot-reload watcher config
├── .env                        # Local environment variables (gitignored)
├── database.yml                # Pop database connection config
├── go.mod / go.sum             # Go module dependencies
├── package.json                # Node dependencies (Svelte, Vite, Tailwind, DaisyUI)
├── svelte.config.js            # Svelte preprocessor config
└── vite.config.js              # Vite build config (outputs to public/assets/)
```

---

## Prerequisites

### All Platforms

| Tool | Install |
|---|---|
| **Go 1.26+** | https://go.dev/dl/ |
| **Buffalo CLI v0.18+** | `go install github.com/gobuffalo/cli/cmd/buffalo@latest` |
| **Node.js 18+** | https://nodejs.org/ |
| **npm 9+** | Included with Node.js |
| **PostgreSQL 14+** | https://www.postgresql.org/download/ |

### Windows Only
- **Git Bash** or **WSL2** recommended for a Unix-like shell experience.
- All commands below work in **PowerShell** unless noted.

### Docker (Alternative)
- **Docker Desktop** — https://www.docker.com/products/docker-desktop/

---

## Configuration

All configuration is done via environment variables. Copy the example below into a `.env` file in the project root (Buffalo loads it automatically):

```env
# Buffalo environment: development | test | production
GO_ENV=development

# Secret used to sign session cookies — CHANGE THIS in production!
SESSION_SECRET=changeme-use-a-random-secret-in-production

# PostgreSQL connection URLs (optional — defaults shown below are used if not set)
# DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/buffalo_app_development?sslmode=disable
# TEST_DATABASE_URL=postgres://postgres:postgres@127.0.0.1:5432/buffalo_app_test?sslmode=disable
```

### Environment Variable Reference

| Variable | Required | Default | Description |
|---|---|---|---|
| `GO_ENV` | No | `development` | Runtime environment |
| `SESSION_SECRET` | **Yes (production)** | — | Cookie signing secret (min 32 chars) |
| `DATABASE_URL` | No | `postgres://postgres:postgres@127.0.0.1:5432/buffalo_app_development?sslmode=disable` | Primary DB connection |
| `TEST_DATABASE_URL` | No | Same host, `..._test` database | Used during `buffalo test` |
| `ADDR` | No | `127.0.0.1` | Bind address (`0.0.0.0` for Docker) |
| `PORT` | No | `3000` | HTTP port |
| `FORCE_SSL` | No | `false` | Force redirect to HTTPS |

---

## Local Development (Without Docker)

The development workflow uses **two terminals** — one for the Go backend, one for the Vite frontend dev server.

### Linux & macOS

**1. Install Go dependencies**
```bash
go mod download
```

**2. Install Node dependencies**
```bash
npm install
```

**3. Start PostgreSQL** (if not already running)
```bash
# macOS with Homebrew
brew services start postgresql@14

# Linux (systemd)
sudo systemctl start postgresql
```

**4. Create the databases**
```bash
buffalo pop create -a
```

**5. Run database migrations**
```bash
buffalo pop migrate
```

**6. Start the backend** (Terminal 1)
```bash
buffalo dev
```

**7. Start the frontend dev server** (Terminal 2)
```bash
npm run dev
```

Open **http://localhost:3000** in your browser. Vite runs on **http://localhost:3001** and serves frontend assets with hot module replacement.

> [!NOTE]
> **Default Seeded Admin User:**
> The database migrations automatically seed a default Super Administrator account:
> - **Email / Username:** `admin@example.com`
> - **Password:** `password123`
> - **Default Organization:** `Default Organization`

---

### Windows

**1. Install Go dependencies**
```powershell
go mod download
```

**2. Install Node dependencies**
```powershell
npm install
```

**3. Start PostgreSQL**

If installed via the PostgreSQL Windows installer, start it from **Services** (`services.msc`) or:
```powershell
# Start via pg_ctl (adjust path to your Postgres installation)
& "C:\Program Files\PostgreSQL\14\bin\pg_ctl.exe" start -D "C:\Program Files\PostgreSQL\14\data"
```

**4. Create the databases**
```powershell
buffalo pop create -a
```

**5. Run database migrations**
```powershell
buffalo pop migrate
```

**6. Start the backend** (Terminal 1)
```powershell
buffalo dev
```

**7. Start the frontend dev server** (Terminal 2)
```powershell
npm run dev
```

Open **http://localhost:3000** in your browser.

> **Tip:** Windows Defender or antivirus may slow down the Go file watcher. Adding the project folder to the exclusion list improves rebuild speed significantly.

---

## Running with Docker

The included `Dockerfile` performs a two-stage build: it compiles the Go binary (with assets embedded) in a Buffalo builder image and produces a minimal Alpine-based production image.

> **Note:** You must build the frontend assets *before* `docker build`, as `buffalo build --static` embeds `public/assets/` into the binary.

### 1. Build frontend assets

```bash
npm install
npm run build
```

### 2. Build the Docker image

```bash
docker build -t todoflow:latest .
```

### 3. Run the container

```bash
docker run -d \
  --name todoflow \
  -p 3000:3000 \
  -e GO_ENV=production \
  -e SESSION_SECRET=your-super-secret-here \
  -e DATABASE_URL="postgres://user:password@host:5432/buffalo_app_production?sslmode=disable" \
  todoflow:latest
```

Open **http://localhost:3000**.

### 4. Run migrations inside the container

```bash
docker exec todoflow /bin/app migrate
```

### Docker Compose (recommended for local Docker dev)

Create a `docker-compose.yml` in the project root:

```yaml
version: '3.8'

services:
  db:
    image: postgres:14-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      GO_ENV: production
      ADDR: 0.0.0.0
      SESSION_SECRET: changeme-replace-in-production
      DATABASE_URL: postgres://postgres:postgres@db:5432/buffalo_app_production?sslmode=disable
    depends_on:
      - db

volumes:
  pgdata:
```

Then run:
```bash
npm run build                     # build frontend assets first
docker compose up --build         # start app + database
docker compose exec app /bin/app migrate  # run migrations
```

---

## Database Management

```bash
# Create all databases defined in database.yml
buffalo pop create -a

# Drop all databases
buffalo pop drop -a

# Run pending migrations
buffalo pop migrate

# Rollback the last migration
buffalo pop migrate down

# Generate a new migration (Fizz format)
buffalo pop generate fizz <migration_name>

# Reset: drop, create, migrate
buffalo pop reset
```

---

## Building for Production

```bash
# 1. Build frontend assets
npm run build

# 2. Build the Go binary (embeds public/assets/ into the binary)
buffalo build -o bin/app

# 3. Run the binary
GO_ENV="production" SESSION_SECRET="<secret>" DATABASE_URL="<url>" ./bin/app
```

> On **Windows**, the binary will be `bin\app.exe`. Run it from PowerShell:
> ```powershell
> $env:GO_ENV="production"; $env:SESSION_SECRET="<secret>"; $env:DATABASE_URL="<url>"; .\bin\app.exe
> ```

### Running Database Migrations in Production

When compiled, the application binary includes all database migrations embedded directly into the executable. You do not need the physical migration source files on disk in production.

To run migrations in your production environment, pass the `migrate` argument to the compiled binary:

```bash
# On Linux/macOS
./bin/app migrate

# On Windows (PowerShell)
.\bin\app.exe migrate
```

To roll back the last applied migration:

```bash
# On Linux/macOS
./bin/app migrate down

# On Windows (PowerShell)
.\bin\app.exe migrate down
```

> [!NOTE]
> Avoid using developer-facing subcommands like `app pop migrate` or `app db migrate` in production. Those commands are intended for local development and expect your physical database configuration and source code files to exist on disk in the current working directory. The standalone `migrate` and `migrate down` commands are optimized for production and run seamlessly from any directory.

### Optimizing CSS Bundle Size (DaisyUI Themes)

By default, the application is configured to import all DaisyUI themes, which adds unused styling rules and theme variables to the compiled production CSS file. 

To optimize and significantly reduce your CSS bundle size for production, you can restrict DaisyUI to load only the themes your application actually uses (e.g., `light` and `dark`):

1. Open [assets/css/application.css](assets/css/application.css).
2. Modify the `@plugin "daisyui"` block to list only your desired theme(s) instead of `all`:

```css
/* Before: Imports all 30+ DaisyUI themes */
@plugin "daisyui" {
  themes: all;
}

/* After: Imports only your selected themes (e.g. light and dark) */
@plugin "daisyui" {
  themes: light, dark;
}
```

3. Rebuild your frontend assets to compile the optimized CSS:
```bash
npm run build
```

---

## Deploying to Google Cloud Run

This application is configured for continuous deployment using Google Cloud Build and Cloud Run, backed by Cloud SQL (PostgreSQL).

### Prerequisites

1.  **Google Cloud Project:** You need an active GCP project with billing enabled.
2.  **Enable APIs:** Enable the following APIs in your project:
    *   Cloud Run API (`run.googleapis.com`)
    *   Cloud Build API (`cloudbuild.googleapis.com`)
    *   Secret Manager API (`secretmanager.googleapis.com`)
    *   Cloud SQL Admin API (`sqladmin.googleapis.com`)
3.  **Cloud SQL Instance:** Create a PostgreSQL instance in Cloud SQL and a database for the application. Make sure the Cloud Run service account has the **Cloud SQL Client** role.

### Secrets Configuration

Use Google Secret Manager to securely store your environment variables:

1.  **Create Secrets:**
    *   Create a secret named `TODOFLOW_SESSION_SECRET` (e.g., a random 32+ character string).
    *   Create a secret named `TODOFLOW_DATABASE_URL` with your Cloud SQL connection string.
        *   Format: `postgres://<user>:<password>@<host>:5432/<database>?sslmode=disable`
        *   *Note: If connecting via Cloud SQL Auth Proxy or Unix sockets, adjust the URL accordingly (e.g., `postgres://<user>:<password>@/buffalo_app_production?host=/cloudsql/<PROJECT_ID>:<REGION>:<INSTANCE_NAME>`).*
2.  **Grant Access:** Ensure the compute service account (used by Cloud Run) has the **Secret Manager Secret Accessor** role for both secrets.

### CI/CD Deployment with Cloud Build

The repository includes a `cloudbuild.yaml` file configured to:
1.  **Run Migrations:** Deploy and execute a temporary Cloud Run Job from source to run the database migrations (`/bin/app migrate`).
2.  **Deploy Application:** Deploy the main Cloud Run Service from source.

Both steps use `gcloud run deploy --source .`, which securely builds the container using Cloud Build (utilizing the multi-stage `Dockerfile`) and deploys it automatically.

**To trigger a manual deployment using Cloud Build:**

```bash
gcloud builds submit --config cloudbuild.yaml .
```

*You can also set up a Cloud Build Trigger to automatically deploy when you push to the `main` branch.*

---

## Stopping the Application

### Local development
- **Backend** (`buffalo dev`): Press `Ctrl+C` in Terminal 1.
- **Frontend** (`npm run dev`): Press `Ctrl+C` in Terminal 2.

### Docker
```bash
# Stop and remove the container
docker stop todoflow && docker rm todoflow

# Docker Compose
docker compose down

# Docker Compose — also remove the database volume (⚠️ destroys data)
docker compose down -v
```

---

## Running Tests

To run the Go automated test suite successfully, the environment must be set to `test` when launching the tests.

First, ensure the test database has been created:
```bash
# Create all databases defined in database.yml (including test)
buffalo pop create -a
```

### Method 1: Using the Buffalo CLI (macOS / Linux)

On Unix-based systems, the Buffalo CLI automatically configures `GO_ENV=test` under the hood:

```bash
# Run all Go tests
buffalo test

# Run tests for a specific package
buffalo test ./actions/...
```

---

### Method 2: Using Standard Go Tools (Windows / IDEs / CI)

On **Windows** (due to a path-escaping bug in Pop's schema loader) or when running tests inside an **IDE** (such as VS Code, GoLand, or Cursor), you must migrate the test database manually once and run the tests with standard Go tools, explicitly setting the environment variable:

#### 1. Migrate the Test Database
```bash
# Run migrations on the test database
buffalo pop migrate -e test
```

#### 2. Run the Suite with `GO_ENV=test`

* **Linux / macOS / Git Bash:**
  ```bash
  GO_ENV=test go test ./...
  ```
* **Windows (PowerShell):**
  ```powershell
  $env:GO_ENV="test"; go test ./...
  ```
* **Windows (Command Prompt):**
  ```cmd
  set GO_ENV=test && go test ./...
  ```

---

## Administrative Tasks (Grifts)

The application provides custom Buffalo "grift" tasks for management and administration operations.

### Account Recovery / Password Reset

If a user gets locked out or you need to manually verify an account and change its password (especially during local development), run the `recover_account` task:

```bash
buffalo task recover_account admin@example.com password123
```

This task will:
- Find the user by their email address.
- Reset and rehash their password to the `<new_password>`.
- Set their `AccountVerified` status to `true`.
- Reset their failed login attempt counter.

---

## Contributing

We welcome contributions to this project! When contributing, please follow these guidelines to maintain code quality and consistency:

### Backend (Go / Buffalo)
- **Code Formatting:** Always run `gofmt -w .` to format Go code.
- **Error Handling:** Check all error returns. Never discard errors with `_`.
- **Pop ORM & Migrations:** 
  - Existing database migrations are immutable. Always generate new ones using `buffalo pop generate fizz <name>`.
  - Always explicitly define single or composite database indexes for foreign keys using `add_index` in Fizz migrations.
  - When loading related records, use Pop ORM's `.Eager()` modifier explicitly as relationship loading is disabled by default.

### Frontend (Svelte 5 / Tailwind / DaisyUI)
- **Svelte 5 Runes:** Use Svelte 5 runes (`$state`, `$derived`, `$effect`, `$props`). Do not use the legacy Options API.
- **UI/UX Styling:** 
  - Prioritize DaisyUI components (e.g., buttons, cards, forms, tables) and Tailwind utility classes.
  - Do not use inline `style="..."` attributes. Avoid writing custom CSS unless absolutely necessary.
- **CSRF Protection:** Always include the `<%= authenticity_token %>` in non-GET Plush forms, or pass the token via the `X-CSRF-Token` header for Svelte/JSON client requests.

### Testing
- Run all automated tests locally to verify your changes. Refer to the **Running Tests** section above for environment setup commands.

