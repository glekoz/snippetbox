# Snippetbox

#### A lightweight, secure Go web application for saving, sharing and managing short text snippets (like Pastebin or GitHub Gists).

Snippetbox is a small, production-minded web app written in Go that lets users create, view and manage short text snippets. It’s designed as a clean example project for building web services with Go. It's currently under development, so new features, bug fixes and overall improvement are coming.

---

## Table of Contents

* [Features](#features)
* [Quick Start](#quick-start)

  * [Prerequisites](#prerequisites)
  * [Run with Docker (recommended)](#run-with-docker-recommended)
* [Configuration](#configuration)

  * [Variables](#variables)
  * [Database setup](#database-setup)
  * [TLS / HTTPS](#tls--https)
* [Directory Layout](#directory-layout)
* [Testing](#testing)
* [Deployment](#deployment)
* [Roadmap / TODO](#roadmap--todo)
* [Contributing](#contributing)
* [Acknowledgments](#acknowledgments)
* [License](#license)

---

## Features

* Create, read and list short text **snippets** (title, content, created at).
* User **authentication** (register, sign in using JWT tokens) and session management.
* Access control: only authenticated users can create or manage their snippets (configurable).
* Persistent storage using a relational database (PostgreSQL by default).
* Secure defaults: TLS support, CSRF protection, input sanitization and secure session cookies.
* Simple, responsive UI rendered with Go HTML templates (server-side rendering).
* Production-ready build and packaging via Docker + `docker-compose`.

---

## Quick Start

### Prerequisites

* Go (recommended version `1.20+`, but the project may work with earlier Go 1.x series)
* PostgreSQL or MySQL instance (local or remote)
* Docker & Docker Compose (recommended for a one-command run)

### Run with Docker (recommended)

This repo includes a `Dockerfile` and `docker-compose.yml` so you can bring the app up quickly.

```bash
# clone the repo
git clone https://github.com/glekoz/snippetbox.git
cd snippetbox

# build and run with docker-compose (creates DB + web service)
docker-compose up --build
```

After the services start, the web app is commonly available at `https://localhost:8000` or `http://localhost:8000` depending on how TLS is configured. Adjust ports in `docker-compose.yml` if needed.

## Configuration

Below are recommended configuration keys you should add or adapt to your code base.

### Variables

| Variable         |                                  Purpose | Example                                                   |
| ---------------- | ---------------------------------------: | --------------------------------------------------------- |
| `DB_DSN`         |   Database connection string (MySQL DSN) | `user:pass@tcp(localhost:3306)/snippetbox?parseTime=true` |
| `PORT`           |                      HTTP(S) listen port | `8000`                                                    |
| `SESSION_SECRET` |   Secret key for signing session cookies | `a-very-secret-string`                                    |
| `TLS_CERT`       |            Path to TLS certificate (PEM) | `/etc/ssl/certs/snippetbox.crt`                           |
| `TLS_KEY`        |            Path to TLS private key (PEM) | `/etc/ssl/private/snippetbox.key`                         |
| `LOG_LEVEL`      | Log verbosity (info, debug, warn, error) | `info`                                                    |

> If your repo uses a config file (YAML / JSON / TOML) you can map the same values there and load them at startup.

### Database setup

1. Create a database and a dedicated user:

```sql
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'snippetbox'@'localhost' IDENTIFIED BY 'yourpassword';
GRANT ALL PRIVILEGES ON snippetbox.* TO 'snippetbox'@'localhost';
FLUSH PRIVILEGES;
```

2. Apply schema. Many forks include `schema.sql` or migration files. If not, create a table like:

```sql
CREATE TABLE snippets (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at DATETIME NULL
);
```

### TLS / HTTPS

For development you can generate a self-signed certificate (many repos include a `Makefile` target for this):

```bash
# simple openssl self-signed cert for local dev
openssl req -x509 -newkey rsa:4096 -nodes -keyout key.pem -out cert.pem -days 365 -subj "/CN=localhost"
export TLS_CERT="$(pwd)/cert.pem"
export TLS_KEY="$(pwd)/key.pem"
```

For production use a real certificate (Let's Encrypt or another CA) and ensure ports and firewall rules allow inbound HTTPS.

---

## Directory Layout (convention)

A typical layout for the project follows this structure (many forks follow this pattern):

```
/ (repo root)
├─ cmd/web/                # main web server package
├─ internal/               # application code not intended for external import
│  ├─ jwtAuth/             # authentication based on JWT tokens
│  ├─ models/              # database models & persistence
│  └─ validator/           # user input validation
├─ ui/                     # static assets + templates
├─ Dockerfile
├─ docker-compose.yml
└─ go.mod
```

---

## Testing

* Run unit tests with `go test ./...`.
* Use table-driven tests and dependency injection for easy testability of handlers and database code.

Examples:

```bash
# run tests
go test ./... -v

# run a single package
go test ./internal/models -run TestInsertSnippet -v
```

---

## Deployment

Suggested deployment approaches:

* **Docker / docker-compose**: simplest for small deployments.
* **Kubernetes**: containerize the app and deploy with a Deployment + Service + Ingress (TLS via cert-manager).

Production checklist:

* Use a managed PostgreSQL instance or a properly secured DB host.
* Use a real TLS certificate (Let's Encrypt / CA).
* Set up backups for your database.
* Configure proper logging and monitoring (stdout logging for Docker or a logging sidecar).
* Run the app as an unprivileged user and follow the principle of least privilege for DB credentials.

---

## Roadmap / TODO

* Add full-text search and tagging for snippets.
* Improve UI/UX: syntax highlighting for code snippets, editor enhancements.
* Add user roles (admin, moderator) and rate limiting to prevent abuse.
* Add automated DB migrations and versioning.
* Implement optional OAuth2 social logins (GitHub, Google).

---

## Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository and create a topic branch (`feature/your-feature`).
2. Open a Pull Request with a clear description and tests.
3. Ensure `go vet` and `go test ./...` pass.
4. Keep changes small and focused.

Consider opening an issue first to discuss larger features.

---

## Acknowledgments

* This project follows the well-known **"Snippetbox"** example from Alex Edwards’ *Let's Go* book and many community forks and adaptations. It’s a great learning resource and served as inspiration for many Go web app examples.
* Thanks to the Go community and authors of helpful libraries and tools that make building, testing and deploying Go web apps straightforward.

---
