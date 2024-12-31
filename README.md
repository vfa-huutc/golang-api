# Golang CMS

This is a Golang project designed to handle a simple web service with user management, roles, permissions, and refresh tokens. It uses MySQL for the database, Docker for containerization, and includes support for migrations, seeding, and authentication.

## Project Structure

The project follows a clean architecture and is organized into the following directories:

```
├── Dockerfile                        # Docker configuration for the application
├── README.md                         # Project documentation
├── cmd                               # Command-line interfaces (CLI)
│   ├── seeder                        # Seeder for initial data population
│   │   └── seeder.go
│   └── server                        # Main entry point for the web server
│       └── main.go
├── configs                           # Configuration files for database, environment variables, JWT, etc.
│   ├── database.go
│   ├── env.go
│   └── jwt.go
├── constants                         # Constants and error handling
│   ├── errors.go
│   └── keys.go
├── docker-compose.yml                # Docker Compose configuration for the app and MySQL
├── docs                              # API documentation
│   └── api_spec.md
├── go.mod                            # Go module dependencies
├── go.sum                            # Go module checksums
├── internal                          # Core application logic
│   ├── database                      # Database migrations and seeding
│   ├── handlers                      # HTTP request handlers
│   ├── middlewares                   # Middlewares for authentication and logging
│   ├── models                        # Data models for the application
│   ├── repositories                  # Repositories for database access
│   ├── routes                        # Routes and routing logic
│   ├── services                      # Business logic for authentication, user, etc.
│   └── utils                         # Utility functions (e.g., for encryption, validation)
├── mysql                             # MySQL database files (e.g., data and logs)
│   └── db/data
├── pkg                               # External packages
│   ├── logger                        # Logger utility
│   └── mailer                        # Mailer for sending emails
├── tests                             # Unit and integration tests
│   └── internal/utils
│       └── security_test.go
```

## Prerequisites

Before getting started, ensure that you have the following installed:

- [Go](https://golang.org/dl/) (Go 1.23 or later)
- [Docker](https://www.docker.com/products/docker-desktop)
- [Docker Compose](https://docs.docker.com/compose/)
- [MySQL](https://www.mysql.com/)

## Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/yourusername/yourproject.git
cd yourproject
```

### 2. Build and run the application using Docker

You can use Docker Compose to set up both the app and the MySQL database:

```bash
docker-compose up --build
```

This will:

- Build the Docker images.
- Start a MySQL container.
- Start the application container.

### 3. Database Migrations

The project includes migrations for creating the necessary tables in the MySQL database.
To apply the migrations:

```bash
migrate -path ./internal/database/migrations -database "mysql://root:root@tcp(127.0.0.1:3306)/golang_db_2" up
```

This will run the migration scripts and populate the database.

### 4. Seeding the Database

To seed the database with initial data (e.g., default users, roles, permissions), run:

```bash
docker-compose exec app go run cmd/seeder/seeder.go
```

### 5. Running the Server

To run the application server:

```bash
air
```

The server will start and be available at `http://localhost:8080`.

## Environment Variables

The following environment variables are required for the application:

- `DB_USERNAME` - MySQL database username.
- `DB_PASSWORD` - MySQL database password.
- `DB_DATABSE`  - MySQL Name of database.
- `DB_PORT`     - MySQL port number of database.

These can be set in the `.env` file or passed directly as environment variables.


Check the `docs/api_spec.md` for a detailed API specification.

## Testing

Run unit tests with the following command:

```bash
go test ./...
```

For specific tests, use:

```bash
go test -v path/to/test
```

### Unit Tests Directory

The test files are located under the `tests` directory. The tests follow the Go testing conventions.

## Contribution Guidelines

1. Fork the repository.
2. Create a feature branch (`git checkout -b feature/feature-name`).
3. Commit your changes (`git commit -am 'Add feature'`).
4. Push to the branch (`git push origin feature/feature-name`).
5. Open a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Key Sections:
1. **Project Structure**: A breakdown of the directories and files with brief descriptions.
2. **Setup Instructions**: Instructions for setting up the project locally, including dependencies and Docker setup.
3. **Environment Variables**: Key environment variables needed for the project to run properly.
4. **Testing**: How to run unit tests in the project.
5. **Contribution Guidelines**: Instructions for contributing to the project.

Feel free to customize this `README.md` based on your actual project requirements.
