package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/mailer"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type config struct {
	port    int
	db_conn string
	smtp    struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
type App struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   config
	db       *database.Queries
	mailer   mailer.Mailer
	wg       sync.WaitGroup
}

func setupConfig() (config, error) {
	godotenv.Load()

	portStr := os.Getenv("server_port")
	port, _ := strconv.Atoi(portStr)
	dbUrl := os.Getenv("db_conn")
	if dbUrl == "" {
		return config{}, errors.New("db_conn was not found in the environment variables")
	}
	smtpHost := os.Getenv("smtp_host")
	smtpPortStr := os.Getenv("smtp_port")
	smtpPort, _ := strconv.Atoi(smtpPortStr)
	smtpUsername := os.Getenv("smtp_username")
	smtpPassword := os.Getenv("smtp_password")
	smtpSender := os.Getenv("smtp_sender")
	cfg := config{
		port:    port,
		db_conn: dbUrl,
		smtp: struct {
			host     string
			port     int
			username string
			password string
			sender   string
		}{
			host:     smtpHost,
			port:     smtpPort,
			username: smtpUsername,
			password: smtpPassword,
			sender:   smtpSender,
		},
	}
	return cfg, nil
}
func main() {
	erorrLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := setupConfig()
	if err != nil {
		erorrLog.Fatal(err)
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.db_conn)
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
		mailer:   mailer.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	if err = app.serve(); err != nil {
		app.errorLog.Fatal(err)
	}
}

func (app *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/v1/healthcheck", app.healthCheckHandler)
	r.Post("/api/v1/users", app.createUserHandler)
	r.Put("/api/v1/users/activated", app.activateUserHandler)
	r.Post("/api/v1/feeds", app.createFeedHandler)
	r.Post("/api/v1/tokens/auth", app.authenticateUserHandler)

	return app.recoverPanic(app.authenticate(r))
}
