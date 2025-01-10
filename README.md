# Go Greenlight
Go Greenlight is a clone of the Open Movie Database core API written in Go. 

The project serves as training for building a REST API app with Go, following (or at least trying to follow) best practices for project structure, dependency management, and code organization.



## Features

- TBD



## Dependencies

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
| PUT    | /v1/activated             | Activate a specific user                        |
| PUT    | /v1/users/password        | Update the password for a specific user         |
| POST   | /v1/tokens/authentication | Generate a new authentication token             |
| POST   | /v1/tokens/password-reset | Generate a new password-reset token             |
| GET    | /debug/vars               | Display application metrics                     |



## Acknowledgements

- This project is strongly based on the Let's Go Further 1.23 book's project, made by Alex Edwards, one of the most prominent Go developers in the community. [Here](https://lets-go-further.alexedwards.net) you can buy it!