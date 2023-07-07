// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package database

import (
	"context"
)

type Querier interface {
	CreateFeed(ctx context.Context, arg CreateFeedParams) (CreateFeedRow, error)
	CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) (FeedFollow, error)
	CreateToken(ctx context.Context, arg CreateTokenParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error)
	DeleteAllUserTokens(ctx context.Context, arg DeleteAllUserTokensParams) error
	DeleteFeedFollow(ctx context.Context, arg DeleteFeedFollowParams) error
	GenerateNextFeedsToFetch(ctx context.Context) ([]GenerateNextFeedsToFetchRow, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
	GetUserByToken(ctx context.Context, arg GetUserByTokenParams) (GetUserByTokenRow, error)
	ListFeedFollow(ctx context.Context) ([]FeedFollow, error)
	ListFeeds(ctx context.Context) ([]ListFeedsRow, error)
	MarkFeedFetched(ctx context.Context, arg MarkFeedFetchedParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) (int32, error)
}

var _ Querier = (*Queries)(nil)