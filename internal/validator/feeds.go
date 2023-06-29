package validator

import (
	"net/url"
	"unicode/utf8"

	"github.com/StelIify/feedbland/internal/database"
)

func ValidateUrl(v *Validator, providedUrl string) {
	v.Check(providedUrl != "", "url", "must be provided")
	_, err := url.ParseRequestURI(providedUrl)
	if err != nil {
		v.Check(false, "url", "must be a valid url")
	}
}

func ValidateFeed(v *Validator, feed *database.Feed) {
	v.Check(feed.Name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(feed.Name) <= 250, "name", "must be not more than 250 characters")

	ValidateUrl(v, feed.Url)
}
