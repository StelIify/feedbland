// Code generated by MockGen. DO NOT EDIT.
// Source: internal\database\querier.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	database "github.com/StelIify/feedbland/internal/database"
	gomock "github.com/golang/mock/gomock"
)

// MockQuerier is a mock of Querier interface.
type MockQuerier struct {
	ctrl     *gomock.Controller
	recorder *MockQuerierMockRecorder
}

// MockQuerierMockRecorder is the mock recorder for MockQuerier.
type MockQuerierMockRecorder struct {
	mock *MockQuerier
}

// NewMockQuerier creates a new mock instance.
func NewMockQuerier(ctrl *gomock.Controller) *MockQuerier {
	mock := &MockQuerier{ctrl: ctrl}
	mock.recorder = &MockQuerierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuerier) EXPECT() *MockQuerierMockRecorder {
	return m.recorder
}

// CreateFeed mocks base method.
func (m *MockQuerier) CreateFeed(ctx context.Context, arg database.CreateFeedParams) (database.CreateFeedRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeed", ctx, arg)
	ret0, _ := ret[0].(database.CreateFeedRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFeed indicates an expected call of CreateFeed.
func (mr *MockQuerierMockRecorder) CreateFeed(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeed", reflect.TypeOf((*MockQuerier)(nil).CreateFeed), ctx, arg)
}

// CreateFeedFollow mocks base method.
func (m *MockQuerier) CreateFeedFollow(ctx context.Context, arg database.CreateFeedFollowParams) (database.FeedFollow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeedFollow", ctx, arg)
	ret0, _ := ret[0].(database.FeedFollow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFeedFollow indicates an expected call of CreateFeedFollow.
func (mr *MockQuerierMockRecorder) CreateFeedFollow(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeedFollow", reflect.TypeOf((*MockQuerier)(nil).CreateFeedFollow), ctx, arg)
}

// CreatePost mocks base method.
func (m *MockQuerier) CreatePost(ctx context.Context, arg database.CreatePostParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockQuerierMockRecorder) CreatePost(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockQuerier)(nil).CreatePost), ctx, arg)
}

// CreateToken mocks base method.
func (m *MockQuerier) CreateToken(ctx context.Context, arg database.CreateTokenParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockQuerierMockRecorder) CreateToken(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockQuerier)(nil).CreateToken), ctx, arg)
}

// CreateUser mocks base method.
func (m *MockQuerier) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.CreateUserRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, arg)
	ret0, _ := ret[0].(database.CreateUserRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockQuerierMockRecorder) CreateUser(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockQuerier)(nil).CreateUser), ctx, arg)
}

// DeleteAllUserTokens mocks base method.
func (m *MockQuerier) DeleteAllUserTokens(ctx context.Context, arg database.DeleteAllUserTokensParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAllUserTokens", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAllUserTokens indicates an expected call of DeleteAllUserTokens.
func (mr *MockQuerierMockRecorder) DeleteAllUserTokens(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAllUserTokens", reflect.TypeOf((*MockQuerier)(nil).DeleteAllUserTokens), ctx, arg)
}

// DeleteFeedFollow mocks base method.
func (m *MockQuerier) DeleteFeedFollow(ctx context.Context, arg database.DeleteFeedFollowParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFeedFollow", ctx, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFeedFollow indicates an expected call of DeleteFeedFollow.
func (mr *MockQuerierMockRecorder) DeleteFeedFollow(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFeedFollow", reflect.TypeOf((*MockQuerier)(nil).DeleteFeedFollow), ctx, arg)
}

// GenerateNextFeedsToFetch mocks base method.
func (m *MockQuerier) GenerateNextFeedsToFetch(ctx context.Context, limit int32) ([]database.GenerateNextFeedsToFetchRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateNextFeedsToFetch", ctx, limit)
	ret0, _ := ret[0].([]database.GenerateNextFeedsToFetchRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateNextFeedsToFetch indicates an expected call of GenerateNextFeedsToFetch.
func (mr *MockQuerierMockRecorder) GenerateNextFeedsToFetch(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateNextFeedsToFetch", reflect.TypeOf((*MockQuerier)(nil).GenerateNextFeedsToFetch), ctx, limit)
}

// GetPostsFollowedByUser mocks base method.
func (m *MockQuerier) GetPostsFollowedByUser(ctx context.Context, userID int64) ([]database.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostsFollowedByUser", ctx, userID)
	ret0, _ := ret[0].([]database.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPostsFollowedByUser indicates an expected call of GetPostsFollowedByUser.
func (mr *MockQuerierMockRecorder) GetPostsFollowedByUser(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostsFollowedByUser", reflect.TypeOf((*MockQuerier)(nil).GetPostsFollowedByUser), ctx, userID)
}

// GetUserByEmail mocks base method.
func (m *MockQuerier) GetUserByEmail(ctx context.Context, email string) (database.GetUserByEmailRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(database.GetUserByEmailRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockQuerierMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockQuerier)(nil).GetUserByEmail), ctx, email)
}

// GetUserByToken mocks base method.
func (m *MockQuerier) GetUserByToken(ctx context.Context, arg database.GetUserByTokenParams) (database.GetUserByTokenRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByToken", ctx, arg)
	ret0, _ := ret[0].(database.GetUserByTokenRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByToken indicates an expected call of GetUserByToken.
func (mr *MockQuerierMockRecorder) GetUserByToken(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByToken", reflect.TypeOf((*MockQuerier)(nil).GetUserByToken), ctx, arg)
}

// ListFeedFollow mocks base method.
func (m *MockQuerier) ListFeedFollow(ctx context.Context) ([]database.FeedFollow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeedFollow", ctx)
	ret0, _ := ret[0].([]database.FeedFollow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeedFollow indicates an expected call of ListFeedFollow.
func (mr *MockQuerierMockRecorder) ListFeedFollow(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeedFollow", reflect.TypeOf((*MockQuerier)(nil).ListFeedFollow), ctx)
}

// ListFeeds mocks base method.
func (m *MockQuerier) ListFeeds(ctx context.Context) ([]database.ListFeedsRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeeds", ctx)
	ret0, _ := ret[0].([]database.ListFeedsRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeeds indicates an expected call of ListFeeds.
func (mr *MockQuerierMockRecorder) ListFeeds(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeeds", reflect.TypeOf((*MockQuerier)(nil).ListFeeds), ctx)
}

// MarkFeedFetched mocks base method.
func (m *MockQuerier) MarkFeedFetched(ctx context.Context, id int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkFeedFetched", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkFeedFetched indicates an expected call of MarkFeedFetched.
func (mr *MockQuerierMockRecorder) MarkFeedFetched(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkFeedFetched", reflect.TypeOf((*MockQuerier)(nil).MarkFeedFetched), ctx, id)
}

// UpdateUser mocks base method.
func (m *MockQuerier) UpdateUser(ctx context.Context, arg database.UpdateUserParams) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, arg)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockQuerierMockRecorder) UpdateUser(ctx, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockQuerier)(nil).UpdateUser), ctx, arg)
}
