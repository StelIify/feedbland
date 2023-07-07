package main

import (
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/StelIify/feedbland/internal/auth"
	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

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

	hashedPassword, err := bcrypt.GenerateFromPassword(user.PasswordHash, 12)
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
	token, err := auth.GenerateToken(user.ID, 3*24*time.Hour, auth.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.db.CreateToken(r.Context(), database.CreateTokenParams{
		Hash:   token.Hash,
		UserID: new_user.ID,
		Expiry: token.Expiry,
		Scope:  token.Scope,
	})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.runInBackground(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          new_user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.html", data)
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

func (app *App) createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	feed := &database.Feed{
		Name: input.Name,
		Url:  input.Url,
	}
	v := validator.NewValidator()
	if validator.ValidateFeed(v, feed); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user := app.contextGetUser(r)
	if auth.IsAnonymous(user) {
		app.errorLog.Println("you are not logged in, it's anonymous")
		app.invalidCredentialsResponse(w, r)
		return
	}
	new_feed, err := app.db.CreateFeed(r.Context(), database.CreateFeedParams{
		Name:   input.Name,
		Url:    input.Url,
		UserID: user.ID,
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
	tokenHash := sha256.Sum256([]byte(input.TokenPlainText))
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
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(input.Password))
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := auth.GenerateToken(user.ID, 24*time.Hour, auth.ScopeAuthentication)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.db.CreateToken(r.Context(), database.CreateTokenParams{
		Hash:   token.Hash,
		UserID: user.ID,
		Expiry: token.Expiry,
		Scope:  token.Scope,
	})

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusCreated, envelope{"auth_token": token}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *App) listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := app.db.ListFeeds(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelope{"feeds": feeds}, nil)
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
	//todo: validate input
	user := app.contextGetUser(r)
	app.errorLog.Println(user.ID)
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
	//todo user should be authenticated, redirect
	app.db.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: id,
	})
}
func (app *App) listFeedFollowHandler(w http.ResponseWriter, r *http.Request) {
	feed_follows, err := app.db.ListFeedFollow(r.Context())
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

func (app *App) scrapeFeedHandler(w http.ResponseWriter, r *http.Request) {
	feedsToScrape, err := app.db.GenerateNextFeedsToFetch(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	type Rss struct {
		XMLName xml.Name `xml:"rss"`
		Text    string   `xml:",chardata"`
		Version string   `xml:"version,attr"`
		Atom    string   `xml:"atom,attr"`
		Channel struct {
			Text  string `xml:",chardata"`
			Title string `xml:"title"`
			Link  struct {
				Text string `xml:",chardata"`
				Href string `xml:"href,attr"`
				Rel  string `xml:"rel,attr"`
				Type string `xml:"type,attr"`
			} `xml:"link"`
			Description   string `xml:"description"`
			Generator     string `xml:"generator"`
			Language      string `xml:"language"`
			LastBuildDate string `xml:"lastBuildDate"`
			Item          []struct {
				Text        string `xml:",chardata"`
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				PubDate     string `xml:"pubDate"`
				Guid        string `xml:"guid"`
				Description string `xml:"description"`
			} `xml:"item"`
		} `xml:"channel"`
	}
	feeds := []Rss{}
	client := http.Client{Timeout: time.Second * 20}
	for _, feed := range feedsToScrape {
		req, err := http.NewRequest("GET", feed.Url, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		response, err := client.Do(req)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		defer response.Body.Close()

		var rssfeed Rss

		err = xml.NewDecoder(response.Body).Decode(&rssfeed)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		feeds = append(feeds, rssfeed)
	}

	app.writeJson(w, http.StatusOK, envelope{"feeds": feeds}, nil)
}
