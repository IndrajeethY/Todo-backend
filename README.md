# Todo-backend

A robust backend API service for managing todos and authentication, implemented in Go. This project provides a scalable, RESTful backend suitable for todo list applications, including user authentication, todo management, and integration with automation bots.

## Features

- **RESTful API** for managing todos (CRUD operations)
- **JWT-based authentication** for user sessions
- **CORS middleware** for secure cross-origin requests
- **SQLite database** support (with configurable path)
- **Configurable environment** via `.env`
- **Bot integration** for automated tasks or notifications
- **Modular codebase** with clear separation of handlers, models, middleware, and routing

## Directory Structure

```
.
├── .env              # Environment variables
├── .gitignore        # Git ignore rules
├── bots/             # Automation bots integration
│   └── bots.go
├── db/               # Database connection and logic
│   └── db.go
├── go.mod            # Go module definition
├── go.sum            # Go dependencies lockfile
├── handlers/         # API endpoint handlers
│   ├── auth.go       # Authentication logic
│   └── todo.go       # Todo management logic
├── main.go           # Application entrypoint
├── middleware/       # Middleware (auth, CORS)
│   ├── auth.go
│   └── cors.go
├── models/           # Data models (users, todos, etc.)
├── router/           # API router setup
│   └── router.go
```

## Getting Started

### Prerequisites

- Go (1.18+ recommended)
- SQLite3

### Setup

1. **Clone the repository**
   ```sh
   git clone https://github.com/IndrajeethY/Todo-backend.git
   cd Todo-backend
   ```

2. **Configure environment variables**

   Edit the `.env` file as needed. Example:
   ```
   PORT=8081
   DATABASE_PATH=./todos.db
   ADMIN_USER=your_admin_username
   ADMIN_PASS=your_admin_password
   ADMIN_ID=your-admin-uuid
   JWT_SECRET=your-long-jwt-secret
   JWT_EXP_HOURS=24
   TG_BOT_TOKEN=your-telegram-bot-token
   OWNER_TG_ID=212121212
   OWNER_DC_ID=546866465655611121213
   DISCORD_TOKEN=your-discord-bot-token
   ```

3. **Install dependencies**
   ```sh
   go mod tidy
   ```

4. **Run the application**
   ```sh
   go run main.go
   ```

   Server will start on the port defined in `.env` (default: 8080).

### Hosting/Deployment

- **Set environment variables** for your deployment host (see `.env` sample above).
- Ensure the database file path (`DATABASE_PATH`) points to a writable location on your host.
- Use a process manager (like `systemd` or `pm2`) or containerize the app with Docker for production reliability.
- Expose the appropriate port and configure TLS as needed for secure access.

## API Overview

- **Authentication**: JWT-based login and registration endpoints.
- **Todos**: Endpoints for creating, listing, updating, and deleting todos.
- **Bots**: Optional integration for automated tasks.

## Code Modules

- `main.go`: Loads environment, initializes database, sets up routing, applies middleware, starts bots, and launches the server.
- `db/`: Handles database connection and migrations.
- `handlers/`: REST API logic for auth and todos.
- `middleware/`: Reusable middleware for CORS and JWT authentication.
- `router/`: Routes API endpoints to handlers.
- `bots/`: Integrates with automation bots.
- `models/`: Defines data structures for users and todos.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

**Repository:** [IndrajeethY/Todo-backend](https://github.com/IndrajeethY/Todo-backend)