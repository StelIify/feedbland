package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type config struct {
	port int
}
type App struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   config
}

func main() {
	erorrLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")

	app := &App{
		config:   cfg,
		errorLog: erorrLog,
		infoLog:  infoLog,
	}

	server := &http.Server{
		Handler:      app.routes(),
		Addr:         fmt.Sprintf(":%d", cfg.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	infoLog.Printf("starting server on :%d", cfg.port)
	erorrLog.Fatal(server.ListenAndServe())
}

func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{"message": "succesful response"}
	err := app.writeJson(w, 200, msg, nil)
	if err != nil {
		app.errorLog.Printf("marshal error: %v", err)
	}
}

func (app *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/v1/healthcheck", app.healthCheckHandler)

	return r
}

func (app *App) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
	response, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for key, header := range headers {
		w.Header()[key] = header
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)

	return nil
}
