# Go Greenlight

Go Greenlight is a movie and user management REST API written in Go. It is a learning project based on the core API of the Open Movie Database, with an emphasis on project structure, dependency management, and code organization.



## Features

- REST endpoints for managing movie and user records
- Authentication, authorization, panic recovery, basic rate limiter and custom metrics middlewares
- User registration and authentication using stateful tokens



## Dependencies

- Go 1.23+
- Golang-migrate
- Make
- Docker



## Third-party packages

- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)
- Postgres driver from [pq](https://github.com/lib/pq)
- [Go-mail](https://github.com/go-mail/mail) to send emails
- Justinas's [Alice](https://github.com/justinas/alice) for a more readable middleware chaining
- [Golang-migrate](https://github.com/golang-migrate/migrate) to manage database migrations
- [Rate](golang.org/x/time/rate) package to implement rate limiters
- The [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) package for password hashing




## Usage

- Clone this repository

  ```bash
  git clone https://github.com/AlessioPani/go-greenlight.git
  ```
  
- Start the Docker services, including PostgreSQL and a test SMTP server:

  ```bash
  docker-compose up --build
  ```

- Apply the database migrations:

  ```bash
  make migrate-up
  ```

- Build and start the application:

  ```bash
  make start
  ```




## API structure

| Method | URL pattern               | Action                                          |
| ------ | ------------------------- | ----------------------------------------------- |
| GET    | /v1/healthcheck           | Show application health and version information |
| GET    | /v1/movies                | List movies                                     |
| POST   | /v1/movies                | Create a new movie                              |
| GET    | /v1/movies/:id            | Show a specific movie                           |
| PATCH  | /v1/movies/:id            | Update a specific movie                         |
| DELETE | /v1/movies/:id            | Delete a specific movie                         |
| POST   | /v1/users                 | Register a new user                             |
| PUT    | /v1/users/activated       | Activate a user                                 |
| PUT    | /v1/users/password        | Update the password for a specific user         |
| POST   | /v1/tokens/authentication | Generate a new authentication token             |
| POST   | /v1/tokens/password-reset | Generate a new password-reset token             |
| POST   | /v1/tokens/activation     | Generate a new activation token                 |
| GET    | /debug/vars               | Display application metrics                     |



## Acknowledgements

- This project is based on the Let's Go Further 1.23 book's project, made by Alex Edwards, one of the most prominent Go developers in the community. [Here](https://lets-go-further.alexedwards.net) you can buy it!
