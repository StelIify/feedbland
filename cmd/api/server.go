package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *App) serve() error {
	server := &http.Server{
		Handler:      app.routes(),
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}
	shutdownError := make(chan error)
	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		signal := <-quit
		app.infoLog.Printf("shutting down server with %s signal", signal.String())
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}
		app.infoLog.Println("completing background tasks")
		app.wg.Wait()
		shutdownError <- nil
	}()

	app.infoLog.Printf("starting server on :%d", app.cfg.port)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdownError
	if err != nil {
		return err
	}

	app.infoLog.Println("stopped server")
	return nil
}
