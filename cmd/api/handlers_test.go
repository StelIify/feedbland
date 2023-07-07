package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/StelIify/feedbland/internal/database"
	"github.com/StelIify/feedbland/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	app := App{}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "api/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}
	app.healthCheckHandler(rr, r)
	result := rr.Result()

	assert.Equal(t, result.StatusCode, http.StatusOK)
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	expectedBody := `{"message":"successful response"}`
	assert.Equal(t, expectedBody, string(body))
}

func TestListFeedsHandler(t *testing.T) {
	expectedFeeds := database.ListFeedsRow{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "TestSubject1",
		Url:       "http://test_url.com",
		UserID:    2,
	}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQuerier := mock.NewMockQuerier(ctrl)
	mockQuerier.EXPECT().ListFeeds(gomock.Any()).Times(1).Return([]database.ListFeedsRow{expectedFeeds}, nil)
	app := App{db: mockQuerier}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/api/v1/feeds", nil)
	assert.NoError(t, err)

	app.routes().ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response database.ListFeedsRow

	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedFeeds.ID, response.ID)
	assert.Equal(t, expectedFeeds.Name, response.Name)
}
