package main

import (
	"log/slog"
	"os"

	"github.com/AlessioPani/go-greenlight/internal/data"
	"github.com/AlessioPani/go-greenlight/internal/data/mocks"
	"github.com/AlessioPani/go-greenlight/internal/mailer"
)

// newTestApplication returns a test application struct.
func newTestApplication() *application {
	app := application{}

	app.config.enableMetrics = false
	app.config.env = "development"
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	app.mailer = mailer.New("localhost", 1025, "", "", "sender")
	app.models = newTestModels()

	return &app
}

// newTestModels returns a test data.Models using all the mocks
// models.
func newTestModels() data.Models {
	return data.Models{
		Movies:      &mocks.MovieModel{},
		Permissions: mocks.PermissionModel{},
		Tokens:      &mocks.TokenModel{},
		Users:       &mocks.UserModel{},
	}
}
