package main

import (
	"log/slog"
	"os"

	"github.com/AlessioPani/go-greenlight/internal/mailer"
)

func newTestApplication() *application {
	var app application

	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	app.mailer = mailer.New("localhost", 1025, "", "", "sender")
	app.config.enableMetrics = false

	return &app
}
