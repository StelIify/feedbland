package validator

import (
	"unicode/utf8"

	"github.com/StelIify/feedbland/internal/database"
)

func ValidateEmail(v *Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(v.Matches(email, EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(utf8.RuneCountInString(password) >= 8, "password", "must be at least 8 characters long")
	v.Check(utf8.RuneCountInString(password) <= 72, "password", "must not be more than 72 characters long")
}

func ValidateUser(v *Validator, user *database.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(user.Name) <= 100, "password", "must not be more than 100 characters long")

	ValidatePassword(v, string(user.PasswordHash))
	ValidateEmail(v, user.Email)
}
