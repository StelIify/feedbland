package main

import (
	"context"
	"net/http"

	"github.com/StelIify/feedbland/internal/database"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *App) contextSetUser(r *http.Request, user *database.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *App) contextGetUser(r *http.Request) *database.User {
	user, ok := r.Context().Value(userContextKey).(*database.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
