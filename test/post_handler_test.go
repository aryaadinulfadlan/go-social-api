package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPostGetDetail(t *testing.T) {
	time.Sleep(time.Second) // TO PREVENT RATE LIMITING
	router := SetupTest()
	postId, _ := uuid.Parse("db5d4c6c-15d5-4555-89ad-555c7f0e4cf9")
	userId, _ := uuid.Parse("e1b4e485-fa48-4d59-8758-e7f988d5cc17")
	tests := []struct {
		name               string
		postId             string
		expectedHttpCode   int
		expectedBodyCode   int
		expectedBodyStatus string
	}{
		{
			name:               "Invalid Post ID",
			postId:             "e3488ac6-7012-4d95-a002-663b9a6f879x",
			expectedHttpCode:   http.StatusBadRequest,
			expectedBodyCode:   http.StatusBadRequest,
			expectedBodyStatus: shared.StatusBadRequest,
		},
		{
			name:               "Post Doesn't Exist",
			postId:             "e3488ac6-7012-4d95-a002-663b9a6f879a",
			expectedHttpCode:   http.StatusNotFound,
			expectedBodyCode:   http.StatusNotFound,
			expectedBodyStatus: shared.StatusNotFound,
		},
		{
			name:               "Valid User",
			postId:             postId.String(),
			expectedHttpCode:   http.StatusOK,
			expectedBodyCode:   http.StatusOK,
			expectedBodyStatus: shared.StatusOK,
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
			body, _ := io.ReadAll(response.Body)
			var responseBody map[string]any
			json.Unmarshal(body, &responseBody)
			assert.Equal(t, test.expectedHttpCode, response.StatusCode)
			assert.Equal(t, test.expectedBodyCode, int(responseBody["code"].(float64)))
			assert.Equal(t, test.expectedBodyStatus, responseBody["status"])
		})
	}
}
