package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostGetDetail(t *testing.T) {
	time.Sleep(time.Second) // TO PREVENT RATE LIMITING
	router := SetupTest()
	postId, _ := uuid.Parse("db5d4c6c-15d5-4555-89ad-555c7f0e4cf9")
	userId, _ := uuid.Parse("e1b4e485-fa48-4d59-8758-e7f988d5cc17")
	tests := []struct {
		name           string
		postId         string
		expectedStatus int
	}{
		{
			name:           "Invalid Post ID",
			postId:         "e3488ac6-7012-4d95-a002-663b9a6f879x",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Post Doesn't Exist",
			postId:         "e3488ac6-7012-4d95-a002-663b9a6f879a",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Valid User",
			postId:         postId.String(),
			expectedStatus: http.StatusOK,
		},
	}
	token, _ := GenerateJWT(userId.String())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/posts/"+test.postId, nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedStatus, response.StatusCode)
		})
	}
}
