package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// version is the application version number.
// TODO: to be generated automatically at build time.
const version = "1.0.0"

// config is a struct that contains the configuration for the application.
type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
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
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	// Initialize a new structured logger.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Open a DB connection pool.
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established successfully")

	// Initialize application config with all the dependencies.
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Get the configured mux with httprouter.
	mux := app.routes()

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

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// Create an empy connection pool with the dsn provided by the config struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Create a context with a 5 seconds deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check the connection pool connection. If it doesn't, within the deadline (5 sec)
	// returns an error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// If no error has been found, return the db connection pool.
	return db, nil
}
