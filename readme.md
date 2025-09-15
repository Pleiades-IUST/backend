# Pleiades Back-end Server

This is the backend project for Pleiades, a mobile network test drive application. The server provides APIs for user authentication, drive data collection, signal processing, and exporting reports, supporting the core functionalities of the Pleiades mobile and web clients.

## Features

- User authentication (signup, login, JWT-based protected routes)
- Drive data collection and storage
- Signal data processing and CSV export
- RESTful API endpoints with CORS support

## Technologies Used

### Go

- Primary backend language
- Uses the Gin framework for HTTP server and routing
- GORM ORM for PostgreSQL database operations
- Environment configuration with `.env` support

### PostgreSQL

- Main relational database for persisting user, drive, and signal data
- Connection and migrations managed in Go via GORM and `golang-migrate`

### Docker

- Dockerfile provided for containerization
- Installs dependencies, sets up database client, and builds the Go application
- Entrypoint script for initializing the database and launching the app

### Additional Tools

- JWT (JSON Web Tokens) for authentication
- Bcrypt for secure password hashing
- CSV export utilities for drive signal data
- CORS middleware for cross-origin API access

## Getting Started

1. Clone the repository and set up your `.env` file with database and server configuration.
2. Build and run with Docker, or run locally with Go.
3. The server exposes RESTful endpoints on port 8080 by default.

## Example Endpoints

- `POST /auth/signup` — Create a new user account
- `POST /auth/login` — Authenticate and receive a JWT
- `GET /drive/all` — Retrieve all drives (user authentication required)
- `GET /drive/csv` — Export drive signal data as CSV

---

For more details, see the code and configuration files in this repository.
