package auth

import "github.com/StelIify/feedbland/internal/database"

var AnonymousUser = &database.User{}

func IsAnymousUser(u *database.User) bool {
	return u == AnonymousUser
}
