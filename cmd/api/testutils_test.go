package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/AlessioPani/go-greenlight/internal/data"
	"github.com/AlessioPani/go-greenlight/internal/data/mocks"
	"github.com/AlessioPani/go-greenlight/internal/mailer"
)

func newTestApplication() *application {
	app := application{}
	wg := sync.WaitGroup{}

	app.config.enableMetrics = false
	app.config.limiter.enabled = true
	app.config.limiter.burst = 10
	app.config.limiter.rps = 10
	app.config.env = "development"
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	app.mailer = mailer.New("localhost", 1025, "", "", "sender")
	app.models = newTestModels()
	app.wg = &wg

	return &app
}

func newTestModels() data.Models {
	return data.Models{
		Movies:      &mocks.MovieModel{},
		Permissions: mocks.PermissionModel{},
		Tokens:      &mocks.TokenModel{},
		Users:       &mocks.UserModel{},
	}
}
