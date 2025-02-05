package main

import (
	"log/slog"
	"os"
	"sync"

	"github.com/AlessioPani/go-greenlight/internal/mailer"
)

func newTestApplication() *application {
	app := application{}
	wg := sync.WaitGroup{}

	app.config.enableMetrics = false
	app.config.env = "development"
	app.logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	app.mailer = mailer.New("localhost", 1025, "", "", "sender")
	app.wg = &wg

	return &app
}
