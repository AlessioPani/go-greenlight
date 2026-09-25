# Go Greenlight

Go Greenlight is a movie and user management REST API written in Go. It is a learning project based on the core API of the Open Movie Database, with an emphasis on project structure, dependency management, and code organization.



## Features

- REST endpoints for managing movie and user records
- Authentication, authorization, panic recovery, basic rate limiter and custom metrics middlewares
- User registration and authentication using stateful tokens



## Dependencies

- Go 1.25+
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
  
- Build and start PostgreSQL, MailHog, and the Go API:

  ```bash
  docker-compose up --build -d
  ```

Compose runs the database migrations after PostgreSQL is ready and starts the API only after the migrations finish. The API is available at `http://localhost:4001` by default; MailHog's inbox is at `http://localhost:8025`. Set `API_PORT` to choose another host port (for example, `API_PORT=4000 docker-compose up --build -d`). To stop the services, run `docker-compose down`. The PostgreSQL data is stored in `db-data/postgres/` and is retained when the containers stop. `make migrate-up` remains available when applying migrations manually to a database.

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


## Register, activate, and authenticate

The examples below use the API started with Docker Compose. Replace the sample password and email with your own values.

### 1. Register an account

```bash
  curl -i -X POST http://localhost:4001/v1/users \
  -H 'Content-Type: application/json' \
  -d '{"name":"John Doe","email":"mail@@example.com","password":"my-password"}'
```

A successful request returns `201 Created`. The API sends an activation email to MailHog. Open `http://localhost:8025` and copy the token from the message. Activation tokens expire after three days and can be used once.

### 2. Activate the account

Send the activation token in the request body:

```bash
  curl -i -X PUT http://localhost:4001/v1/users/activated \
  -H 'Content-Type: application/json' \
  -d '{"token":"<activation_token_from_mailhog>"}'
```

The response includes the activated user. Registration grants the `movies:read` permission by default.

### 3. Create an authentication token

```bash
  curl -i -X POST http://localhost:4001/v1/tokens/authentication \
  -H 'Content-Type: application/json' \
  -d '{"email":"ada@example.com","password":"correct-horse-battery"}'
```

Copy the `authentication_token.token` value from the response. Authentication tokens expire after 24 hours.

### 4. Call an authenticated endpoint

Pass the token as a Bearer token. For example, listing movies requires `movies:read`:

```bash
  curl -i http://localhost:4001/v1/movies \
  -H 'Authorization: Bearer <authentication_token>'
```

To request a password-reset token, post the account email to `/v1/tokens/password-reset`; the instructions are delivered through MailHog. The password-reset token expires after 45 minutes and is submitted with the new password to `PUT /v1/users/password`.



## Acknowledgements

- This project is based on the Let's Go Further 1.23 book's project, made by Alex Edwards, one of the most prominent Go developers in the community. [Here](https://lets-go-further.alexedwards.net) you can buy it!
