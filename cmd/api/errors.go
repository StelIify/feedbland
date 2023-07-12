package main

import (
	"net/http"
)

type envelope map[string]interface{}

func (app *App) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJson(w, status, env, nil)
	if err != nil {
		app.errorLog.Println(err)
		w.WriteHeader(500)
	}
}

func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
func (app *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorLog.Println(err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *App) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *App) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *App) invalidAuthTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "ivalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *App) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this recourse"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (app *App) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this recourse"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
