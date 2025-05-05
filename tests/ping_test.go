package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	router := SetupTest()
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{name: "OK", expectedStatus: http.StatusOK},
		{name: "OK", expectedStatus: http.StatusOK},
		{name: "OK", expectedStatus: http.StatusOK},
		{name: "OK", expectedStatus: http.StatusOK},
		{name: "OK", expectedStatus: http.StatusOK},
		{name: "Too Many Request", expectedStatus: http.StatusTooManyRequests},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedStatus, response.StatusCode)
		})
	}
}

func TestBasicAuthentication(t *testing.T) {
	router := SetupTest()
	tests := []struct {
		name           string
		base64         string
		expectedStatus int
	}{
		{name: "Valid Base64", base64: "YWRpbnVsOmFkaW51bDEyMw==", expectedStatus: http.StatusOK},
		{name: "Invalid Base64", base64: "Invalid Base64", expectedStatus: http.StatusUnauthorized},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/basic", nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Basic %s", test.base64))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedStatus, response.StatusCode)
		})
	}
}

func TestBearerAuthentication(t *testing.T) {
	router := SetupTest()
	user_id, _ := uuid.Parse("50b466de-2de4-4e40-bdec-08270f23a8c8")
	token, _ := GenerateJWT(user_id.String())
	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{name: "Valid Token", token: token, expectedStatus: http.StatusOK},
		{name: "Invalid Token", token: "Invalid Token", expectedStatus: http.StatusUnauthorized},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/bearer", nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.token))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedStatus, response.StatusCode)
		})
	}
}
