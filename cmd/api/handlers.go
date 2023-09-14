package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/StelIify/feedbland/internal/auth"
	"github.com/StelIify/feedbland/internal/data"
	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/validator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// @Summary Health Check
// @Description Perform a health check on the API
// @ID healthCheck
// @Tags Health
// @Produce json
// @Success 200 {object} string
// @Router /api/v1/healthcheck [get]
func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{"message": "successful response"}
	err := app.writeJson(w, 200, msg, nil)
	if err != nil {
		app.errorLog.Printf("marshal error: %v", err)
	}
}

func (app *App) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := &database.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: []byte(input.Password),
	}
	v := validator.NewValidator()
	if validator.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	hashedPassword, err := auth.GenerePasswordHash(input.Password, 12)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	new_user, err := app.db.CreateUser(r.Context(), database.CreateUserParams{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		var pg_err *pgconn.PgError
		if errors.As(err, &pg_err) && pg_err.Code == pgerrcode.UniqueViolation {
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	token, err := auth.CreateUserToken(r.Context(), app.db, new_user.ID, 3*24*time.Hour, auth.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.runInBackground(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          new_user.ID,
		}
		err = app.mailer.Send(new_user.Email, "user_welcome.html", data)
		if err != nil {
			app.errorLog.Printf("problem during the process of sending welcome email: %v", err)
		}
	})
	err = app.writeJson(w, http.StatusAccepted, envelope{"user": new_user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlainText string `json:"token"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.NewValidator()
	if validator.ValidateTokenPlainText(v, input.TokenPlainText); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	tokenHash := auth.GenerateTokenHash(input.TokenPlainText)
	user, err := app.db.GetUserByToken(r.Context(), database.GetUserByTokenParams{
		Hash:   tokenHash[:],
		Scope:  auth.ScopeActivation,
		Expiry: time.Now(),
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	userVersion, err := app.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Activated:    true,
		ID:           user.ID,
		Version:      user.Version,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	app.writeJson(w, http.StatusOK, userVersion, nil)
}

func (app *App) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.NewValidator()
	validator.ValidateEmail(v, input.Email)
	validator.ValidatePassword(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.db.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err, _ = auth.ValidateCredentials(user.PasswordHash, input.Password)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := auth.CreateUserToken(r.Context(), app.db, user.ID, 24*time.Hour, auth.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// cookie := http.Cookie{
	// 	Name:     "auth_token",
	// 	Value:    token.Plaintext,
	// 	Expires:  token.Expiry,
	// 	HttpOnly: true,
	// }
	// http.SetCookie(w, &cookie)
	err = app.writeJson(w, http.StatusOK, envelope{"auth_token": token}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// @Summary List Feeds
// @Description Get a list of feeds
// @ID listFeeds
// @Tags Feeds
// @Produce json
// @Success 200 {object} database.ListFeedsRow
// @Failure 500 {object} string
// @Router /api/v1/feeds [get]
func (app *App) listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.NewValidator()
	qs := r.URL.Query()
	filters := data.NewFeedsFilters(qs, v)

	if validator.ValidateFilters(v, filters.Offset, filters.Limit, filters.Sort, filters.SortSafelist); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	feeds, err := app.customQueries.ListAllFeeds(r.Context(), filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	feedCound, err := app.db.CountFeeds(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	metadata := data.NewMetadata(feedCound, len(feeds), r.URL.Path, filters)
	err = app.writeJson(w, http.StatusOK, envelope{"feeds": feeds, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// TODO refactor handler
func (app *App) createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Url string `json:"url"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	rssFeed, err := UrlToFeed(input.Url)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if rssFeed.Channel.Image.URL == "" {
		// rss feed does not contain image informatin
		// we need to scrape it
		app.errorLog.Println("rss feed does not have image information")
		return
	}
	//download image
	response, err := http.Get(rssFeed.Channel.Image.URL)
	if err != nil {
		fmt.Println("get request error", err)
	}
	defer response.Body.Close()
	imgBytes, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("read response body error", err)
	}
	//save image to s3 bucket, create image from returned url
	u, err := url.Parse(rssFeed.Channel.Image.URL)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	filename := path.Base(u.Path)
	feedsStore := fmt.Sprintf("%s/%s", "feedsImg", filename)
	result, err := app.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("feebland"),
		Key:    aws.String(feedsStore),
		Body:   bytes.NewReader(imgBytes),
		ACL:    "public-read",
	})
	if err != nil {
		app.errorLog.Println(err)
	}
	// create image
	image_id, err := app.db.CreateImage(r.Context(), database.CreateImageParams{
		Url: pgtype.Text{
			String: result.Location,
			Valid:  true,
		},
		Name: rssFeed.Channel.Image.Title,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	feed := &database.Feed{
		Name: rssFeed.Channel.Title,
		Description: pgtype.Text{
			String: rssFeed.Channel.Description,
			Valid:  true,
		},
		ImageID: pgtype.Int8{
			Int64: image_id,
			Valid: true,
		},
		Url: input.Url,
	}
	v := validator.NewValidator()
	if validator.ValidateFeed(v, feed); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user := app.contextGetUser(r)
	new_feed, err := app.db.CreateFeed(r.Context(), database.CreateFeedParams{
		Name: rssFeed.Channel.Title,
		Description: pgtype.Text{
			String: rssFeed.Channel.Description,
			Valid:  true,
		},
		Url:    input.Url,
		UserID: user.ID,
		ImageID: pgtype.Int8{
			Int64: image_id,
			Valid: true,
		},
	})

	if err != nil {
		var pg_err *pgconn.PgError
		if errors.As(err, &pg_err) && pg_err.Code == pgerrcode.UniqueViolation {
			v.AddError("url", "feed with this url already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	_, err = app.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: new_feed.ID,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusCreated, envelope{"feed": new_feed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) createFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FeedID int64 `json:"feed_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	feedFollow, err := app.db.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		UserID: user.ID,
		FeedID: input.FeedID,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusCreated, feedFollow, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	app.db.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: id,
	})
}

func (app *App) listFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	feed_follows, err := app.db.ListFeedFollow(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, feed_follows, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) listPostsFollowedByUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	posts, err := app.db.GetPostsFollowedByUser(r.Context(), user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"posts": posts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) listPostsForFeedHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	posts, err := app.db.GetPostsForFeed(r.Context(), id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"posts": posts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *App) listPosts(w http.ResponseWriter, r *http.Request) {
	v := validator.NewValidator()
	qs := r.URL.Query()
	filters := data.NewPostsFilters(qs, v)

	if validator.ValidateFilters(v, filters.Offset, filters.Limit, filters.Sort, filters.SortSafelist); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	posts, err := app.customQueries.ListAllPosts(r.Context(), filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	count, err := app.db.CountPosts(r.Context(), filters.Title)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	metadata := data.NewMetadata(count, len(posts), r.URL.Path, filters)
	err = app.writeJson(w, http.StatusOK, envelope{"posts": posts, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
