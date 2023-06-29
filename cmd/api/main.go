package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/StelIify/feedbland/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type config struct {
	port int
}
type App struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   config
	db       *database.Queries
}

func main() {
	erorrLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "API server port")

	godotenv.Load()

	dbUrl := os.Getenv("db_conn")
	if dbUrl == "" {
		erorrLog.Fatal("db_conn is not found in the enviroment")
	}
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		erorrLog.Fatal("Can't connect to the database", err)
	}
	defer dbpool.Close()

	db := database.New(dbpool)

	app := &App{
		config:   cfg,
		errorLog: erorrLog,
		infoLog:  infoLog,
		db:       db,
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

func (app *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/v1/healthcheck", app.healthCheckHandler)
	r.Post("/api/v1/users", app.createUserHandler)
	r.Post("/api/v1/feeds", app.createFeedHandler)

	return r
}
