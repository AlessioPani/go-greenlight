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

// serve configures and runs an HTTP server
// and listens for shutdown signals in a background goroutine.
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

		// Create a context with a 30-second timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Call Shutdown() on our server to stop accepting any HTTP request,
		// passing in the context we just made.
		// Shutdown() will return nil if the graceful shutdown was successful, or an
		// error (which may happen because the shutdown didn't complete before
		// the 30-second context deadline is hit).
		// We send back this return value to the shutdownError channel.
		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		// Wait for all the background goroutines to finish their task before
		// shutting down the server.
		app.logger.Info("completing background tasks", "addr", server.Addr)

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting server", "addr", server.Addr, "env", app.config.env)

	// Calling Shutdown() will cause ListenAndServer to immediately return an
	// ErrServerClosed error. If this occur, it means that the graceful shutdown
	// has started correctly.
	// Otherwise, return the unexpected error.
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait for the shutdown result and return any error.
	// Otherwise, log that the server stopped and return nil.
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", "addr", server.Addr)

	return nil
}
