package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/StelIify/feedbland/internal/auth"
	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/validator"
	"github.com/jackc/pgx/v5"
)

func (app *App) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// run in case of panic
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})

}

func (app *App) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		autorizationHeader := r.Header.Get("Authorization")
		if autorizationHeader == "" {
			r = app.contextSetUser(r, auth.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(autorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)
			return
		}
		token := headerParts[1]

		v := validator.NewValidator()
		if validator.ValidateTokenPlainText(v, token); !v.Valid() {
			app.invalidAuthTokenResponse(w, r)
			return
		}
		tokenHash := sha256.Sum256([]byte(token))
		user, err := app.db.GetUserByToken(r.Context(), database.GetUserByTokenParams{
			Hash:   tokenHash[:],
			Scope:  auth.ScopeAuthentication,
			Expiry: time.Now(),
		})
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				app.invalidAuthTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}
		userdb := &database.User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    time.Time{},
			Name:         user.Name,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			Activated:    user.Activated,
			Version:      user.Version,
		}
		r = app.contextSetUser(r, userdb)
		next.ServeHTTP(w, r)
	})
}
