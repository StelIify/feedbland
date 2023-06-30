package main

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"time"

	"github.com/StelIify/feedbland/internal/auth"
	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/validator"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{"message": "succesful response"}
	time.Sleep(5 * time.Second)
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

	new_feed, err := app.db.CreateFeed(r.Context(), database.CreateFeedParams{
		Name: input.Name,
		Url:  input.Url,
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
	err = app.writeJson(w, http.StatusCreated, envelope{"feed": new_feed}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
