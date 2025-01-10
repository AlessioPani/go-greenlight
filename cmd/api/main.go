package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// version is the application version number.
// TODO: to be generated automatically at build time.
const version = "1.0.0"

// config is a struct that contains the configuration for the application.
type config struct {
	port int
	env  string
}

// application is a struct that contains the dependencies for the application.
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// Instance of the config struct.
	var cfg config

	// Fill the config struct by parsing command line parameters.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new structured logger.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize application config with all the dependencies.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Create a new mux.
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// Start the server and check for errors.
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", "addr", server.Addr, "env", cfg.env)

	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
