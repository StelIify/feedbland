// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package database

import (
	"time"
)

type Feed struct {
	ID            int64     `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Url           string    `json:"url"`
	UserID        int64     `json:"user_id"`
	LastFetchedAt time.Time `json:"last_fetched_at"`
}

type FeedFollow struct {
	FeedID    int64     `json:"feed_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Post struct {
	ID          int64     `json:"id"`
	FeedID      int64     `json:"feed_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
}

type Token struct {
	Hash   []byte    `json:"hash"`
	UserID int64     `json:"user_id"`
	Expiry time.Time `json:"expiry"`
	Scope  string    `json:"scope"`
}

type User struct {
	ID           int64     `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"password_hash"`
	Activated    bool      `json:"activated"`
	Version      int32     `json:"version"`
}
