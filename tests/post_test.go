package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	router := SetupTest()
	user_id, _ := uuid.Parse("50b466de-2de4-4e40-bdec-08270f23a8c8")
	token, err := GenerateJWT(user_id.String())
	assert.Nil(t, err)
	payload := &db.Post{
		Id:        uuid.New(),
		UserId:    user_id,
		Title:     "Title test",
		Content:   "Content test",
		Tags:      []string{"haha", "hihi"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	jsonData, _ := json.Marshal(payload)
	requestBody := bytes.NewReader(jsonData)
	request := httptest.NewRequest(http.MethodPost, "/v1/posts", requestBody)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, http.StatusCreated, int(responseBody["code"].(float64)))
	assert.Equal(t, "Title test", responseBody["data"].(map[string]any)["title"])
}

// func TestGetDetail(t *testing.T) {
// 	router := SetupTest()
// 	user_id, _ := uuid.Parse("50b466de-2de4-4e40-bdec-08270f23a8c8")
// 	authenticator := auth.NewJWTAuthenticator(config.SecretKey)
// 	exp := time.Now().Add(config.Auth.TokenExp).UTC()
// 	_, err := authenticator.GenerateJWT(user_id.String(), exp)
// 	assert.Nil(t, err)
// 	post_id := uuid.New()
// 	postRepo := post.NewRepository(db.DB)
// 	post := &db.Post{
// 		Id:        post_id,
// 		UserId:    user_id,
// 		Title:     "Title here",
// 		Content:   "Content here again",
// 		Tags:      []string{"haha", "hihi"},
// 		CreatedAt: time.Now(),
// 		UpdatedAt: time.Now(),
// 	}
// 	err = postRepo.Create(ctx, post)
// 	assert.Nil(t, err)
// 	// request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/posts/%s", post_id.String()), nil)
// 	// request.Header.Set("Content-Type", "application/json")
// 	// request.Header.Set("Authorization", "Bearer "+token)
// 	// recorder := httptest.NewRecorder()
// 	// router.ServeHTTP(recorder, request)
// 	// response := recorder.Result()
// 	// assert.Equal(t, http.StatusOK, response.StatusCode)
// }
