package main

import (
	"errors"
	"net/http"

	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/validator"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (app *App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	msg := map[string]string{"message": "succesful response"}
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
	err = app.writeJson(w, http.StatusCreated, envelope{"user": new_user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `"json:name"`
		Url  string `"json:url"`
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
