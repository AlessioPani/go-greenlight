package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// server is a method that configures and runs an http.Server
// and checks for shutdowns in background with a gorouting.
func (app *application) serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	// Starts a background goroutine to check for shutdowns.
	go func() {
		// Intercept syscalls to interrupt or terminate the process.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// Wait to receive a syscall to our quit channel.
		s := <-quit
		app.logger.Info("shutting down server", "signal", s.String())

		// Create a context with a 30 seconds of timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Call Shutdown() on our server to stop accepting any HTTP request,
		// passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because the shutdown didn't complete before
		// the 30-second context deadline is hit).
		// We send back this return value to the shutdownError channel.
		shutdownError <- server.Shutdown(ctx)
	}()

	app.logger.Info("starting server", "addr", server.Addr, "env", app.config.env)

	// Calling Shutdown() will cause ListenAndServer to immediately return an
	// ErrServerClosed error. If this occur, it means that the graceful shutdown
	// has started correctly.
	// Otherwise we return the error because it is a unexpected behaviour.
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait the result of Shutdown. If it is an error, return it back.
	// Otherwise log a "stopped server" message and returns no error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", server.Addr)

	return nil
}
