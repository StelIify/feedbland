package main

import (
	"encoding/json"
	"net/http"

	"github.com/StelIify/feedbland/internal/database"
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
	type input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	var i input
	err := decoder.Decode(&i)
	if err != nil {
		//responde with bad request in json
		app.errorLog.Print(err)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(i.Password), 12)
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	new_user, err := app.db.CreateUser(r.Context(), database.CreateUserParams{
		Name:         i.Name,
		Email:        i.Email,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		app.errorLog.Print(err)
		return
	}
	app.writeJson(w, 200, new_user, nil)
}
