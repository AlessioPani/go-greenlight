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

	"github.com/AlessioPani/go-greenlight/internal/data"
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
		// maxOpenConns limits the number of open connections (active and idle) imposed
		// by both database and infrastructure (Postgres has an hard limit of 100 connections).
		maxOpenConns int
		// maxIdleConns limits the number of idle connections.
		maxIdleConns int
		// maxIdleTime limits the max duration for a connection to be in the idle status.
		// After this period of time, the resource will be freed up.
		maxIdleTime time.Duration
	}
	limiter struct {
		// rps is the number of requests per second.
		rps float64
		// burst limits the maximum number of requests on a single burst.
		burst int
		// enabled is a flag that enable (true) or disable (false) the rate limiter.
		enabled bool
	}
}

// application is a struct that contains the dependencies for the application.
type application struct {
	config config
	logger *slog.Logger
	models data.Models
}

func main() {
	// Instance of the config struct.
	var cfg config

	// Fill the config struct by parsing command line parameters.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
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
		models: data.NewModels(db),
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

	// Configure the database connection pool.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	// Create a context with a 5 seconds deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check the connection pool connection. If it doesn't within the deadline (5 sec)
	// returns an error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// If no error has been found, return the db connection pool.
	return db, nil
}
