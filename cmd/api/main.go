package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/StelIify/feedbland/docs"
	"github.com/StelIify/feedbland/internal/data"
	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/mailer"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type cfg struct {
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
	errorLog      *log.Logger
	infoLog       *log.Logger
	cfg           cfg
	db            database.Querier
	customQueries *data.CustomQueries
	mailer        mailer.Mailer
	uploader      *manager.Uploader
	wg            sync.WaitGroup
}

func setupConfig() (cfg, error) {
	godotenv.Load()

	portStr := os.Getenv("server_port")
	port, _ := strconv.Atoi(portStr)
	dbUrl := os.Getenv("db_conn")
	if dbUrl == "" {
		return cfg{}, errors.New("db_conn was not found in the environment variables")
	}
	smtpHost := os.Getenv("smtp_host")
	smtpPortStr := os.Getenv("smtp_port")
	smtpPort, _ := strconv.Atoi(smtpPortStr)
	smtpUsername := os.Getenv("smtp_username")
	smtpPassword := os.Getenv("smtp_password")
	smtpSender := os.Getenv("smtp_sender")
	cfg := cfg{
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

// @title FeedBland
// @version 1.0
// @description backend for blog aggregator

// @host localhost:8080
// @BasePath /

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
	customQueries := data.NewCustomQueries(dbpool)

	//aws s3 setup
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	client := s3.NewFromConfig(awsConfig)
	uploader := manager.NewUploader(client)

	app := &App{
		cfg:           cfg,
		errorLog:      erorrLog,
		infoLog:       infoLog,
		db:            db,
		customQueries: customQueries,
		uploader:      uploader,
		mailer:        mailer.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}
	// go app.fetchFeedsWorker(10, time.Hour*24)

	if err = app.serve(); err != nil {
		app.errorLog.Fatal(err)
	}
}
func resolveURL(baseURL *url.URL, imageURL string) string {
	imageURL = strings.TrimSpace(imageURL)
	if imageURL == "" {
		return ""
	}

	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return ""
	}

	// Check if the URL is relative
	if parsedURL.Scheme == "" && parsedURL.Host == "" {
		resolvedURL := baseURL.ResolveReference(parsedURL)
		return resolvedURL.String()
	}

	return imageURL
}

func (app *App) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/api/v1/healthcheck", app.healthCheckHandler)

	r.Post("/api/v1/users", app.createUserHandler)
	r.Put("/api/v1/users/activated", app.activateUserHandler)
	r.Post("/api/v1/tokens/auth", app.authenticateUserHandler)

	r.Post("/api/v1/feeds", app.requireAuth(app.createFeedHandler))
	r.Get("/api/v1/feeds", app.listFeedsHandler)

	r.Post("/api/v1/feed_follows", app.requireAuth(app.createFeedFollowHandler))
	r.Delete("/api/v1/feed_follows/{id}", app.requireAuth(app.deleteFeedFollowHandler))
	r.Get("/api/v1/feed_follows", app.requireAuth(app.listFeedFollowHandler))

	r.Get("/api/v1/posts", app.requireAuth(app.listPostsFollowedByUserHandler))
	r.Get("/api/v1/feeds/{id}/posts", app.requireAuth(app.listPostsForFeedHandler))
	r.Get("/api/v1/allposts", app.listPosts) //@todo test endpoint, delete later

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	return app.recoverPanic(app.authenticate(r))
}
