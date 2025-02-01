# Go Greenlight
Go Greenlight is a clone of the Open Movie Database core API written in Go. 

The project serves as training for building a REST API app with Go, following (or at least trying to follow) best practices for project structure, dependency management, and code organization.



## Features

- TBD

  

## Dependencies

- Julien Schmidt's [httprouter](https://github.com/julienschmidt/httprouter)

- Postgres driver from [pq](https://github.com/lib/pq)

- [Go-mail](https://github.com/go-mail/mail) to send emails

- Justinas's [Alice](https://github.com/justinas/alice) for a more readable middleware chaining

- [Golang-migrate](https://github.com/golang-migrate/migrate) to manage database migrations

- [Rate](golang.org/x/time/rate) package to implement rate limiters

- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) package for hashing alghoritms

- Make

- Docker

  

## Todo

- Tests

  


## Usage

- TBD  

  


## API structure

| Method | URL pattern               | Action                                          |
| ------ | ------------------------- | ----------------------------------------------- |
| GET    | /v1/healthcheck           | Show application health and version information |
| GET    | /v1/movies                | Show the details of all movies                  |
| POST   | /v1/movies                | Create a new movie                              |
| GET    | /v1/movies/:id            | Show the details of a specific movie            |
| PATCH  | /v1/movies/:id            | Update the details of a specific movie          |
| DELETE | /v1/movies/:id            | Delete a specific movie                         |
| POST   | /v1/users                 | Register a new user                             |
| PUT    | /v1/users/activated       | Activate a specific user                        |
| PUT    | /v1/users/password        | Update the password for a specific user         |
| POST   | /v1/tokens/authentication | Generate a new authentication token             |
| POST   | /v1/tokens/password-reset | Generate a new password-reset token             |
| POST   | /v1/tokens/activation     | Generate a new activation token                 |
| GET    | /debug/vars               | Display application metrics                     |



## Acknowledgements

- This project is based on the Let's Go Further 1.23 book's project, made by Alex Edwards, one of the most prominent Go developers in the community. [Here](https://lets-go-further.alexedwards.net) you can buy it!
