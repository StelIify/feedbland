package data

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/StelIify/feedbland/internal/validator"
)

const baseUrl = "http://localhost:8080"

var (
	defaultTitle  = ""
	defaultLimit  = 20
	defaultOffset = 0

	defaultSortColumnPosts   = "-published_at"
	defaultSortSafeListPosts = []string{"id", "title", "published_at", "-id", "-title", "-published_at"}

	defaultSortColumnFeeds   = "created_at"
	defaultSortSafeListFeeds = []string{"id", "name", "created_at", "-id", "-name", "-created_at"}
)

type Filters struct {
	Title        string
	Limit        int
	Offset       int
	Sort         string
	SortSafelist []string
}

func NewPostsFilters(qs url.Values, v *validator.Validator) Filters {
	return Filters{
		Title:        ReadString(qs, "title", defaultTitle),
		Limit:        ReadInt(qs, "limit", defaultLimit, v),
		Offset:       ReadInt(qs, "offset", defaultOffset, v),
		Sort:         ReadString(qs, "sort", defaultSortColumnPosts),
		SortSafelist: defaultSortSafeListPosts,
	}
}
func NewFeedsFilters(qs url.Values, v *validator.Validator) Filters {
	return Filters{
		Title:        ReadString(qs, "title", defaultTitle),
		Limit:        ReadInt(qs, "limit", defaultLimit, v),
		Offset:       ReadInt(qs, "offset", defaultOffset, v),
		Sort:         ReadString(qs, "sort", defaultSortColumnFeeds),
		SortSafelist: defaultSortSafeListFeeds,
	}
}

func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}

	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func ReadString(qs url.Values, key, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}
	return s
}

func ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be integer value")
		return defaultValue
	}
	return i
}

type Metadata struct {
	Count    int64   `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

func NewMetadata(count int64, posts int, path string, f Filters) *Metadata {
	return &Metadata{
		Count:    count,
		Next:     GetNextUrl(posts, path, f),
		Previous: GetPreviousUrl(path, f),
	}
}

func (f *Filters) NextOffset() int {
	return f.Offset + f.Limit
}
func (f *Filters) PreviousOffset() int {
	return f.Offset - f.Limit
}

func GetNextUrl(posts int, path string, f Filters) *string {
	if posts != f.Limit {
		return nil
	}

	if f.Title == "" {
		nextLink := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl+path, f.Limit, f.NextOffset())
		return &nextLink
	} else {
		nextLink := fmt.Sprintf("%s?title=%s&limit=%d&offset=%d", baseUrl+path, f.Title, f.Limit, f.NextOffset())
		return &nextLink
	}
}

func GetPreviousUrl(path string, f Filters) *string {
	if f.Offset <= 0 {
		return nil
	}
	if f.Title == "" {
		prevLink := fmt.Sprintf("%s?limit=%d&offset=%d", baseUrl+path, f.Limit, f.PreviousOffset())
		return &prevLink
	} else {
		prevLink := fmt.Sprintf("%s?title=%s&limit=%d&offset=%d", baseUrl+path, f.Title, f.Limit, f.PreviousOffset())
		return &prevLink
	}
}
