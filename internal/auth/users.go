package auth

import "github.com/StelIify/feedbland/internal/database"

var AnonymousUser = &database.User{}

func IsAnonymous(u *database.User) bool {
	return u == AnonymousUser
}

func IsActivated(u *database.User) bool {
	return u.Activated
}
