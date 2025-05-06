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

func TestRateLimiter(t *testing.T) {
	router := SetupTest()
	tests := []struct {
		name             string
		expectedHttpCode int
	}{
		{name: "OK", expectedHttpCode: http.StatusOK},
		{name: "OK", expectedHttpCode: http.StatusOK},
		{name: "OK", expectedHttpCode: http.StatusOK},
		{name: "OK", expectedHttpCode: http.StatusOK},
		{name: "OK", expectedHttpCode: http.StatusOK},
		{name: "Too Many Request", expectedHttpCode: http.StatusTooManyRequests},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedHttpCode, response.StatusCode)
		})
	}
}

func TestBasicAuthentication(t *testing.T) {
	time.Sleep(time.Second) // TO PREVENT RATE LIMITING
	router := SetupTest()
	tests := []struct {
		name             string
		base64           string
		expectedHttpCode int
	}{
		{name: "Valid Base64", base64: "YWRpbnVsOmFkaW51bDEyMw==", expectedHttpCode: http.StatusOK},
		{name: "Invalid Base64", base64: "Invalid Base64", expectedHttpCode: http.StatusUnauthorized},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/basic", nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Basic %s", test.base64))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedHttpCode, response.StatusCode)
		})
	}
}

func TestBearerAuthentication(t *testing.T) {
	router := SetupTest()
	user_id, _ := uuid.Parse("e1b4e485-fa48-4d59-8758-e7f988d5cc17")
	token, _ := GenerateJWT(user_id.String())
	tests := []struct {
		name             string
		token            string
		expectedHttpCode int
	}{
		{name: "Valid Token", token: token, expectedHttpCode: http.StatusOK},
		{name: "Invalid Token", token: "Invalid Token", expectedHttpCode: http.StatusUnauthorized},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/bearer", nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", test.token))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedHttpCode, response.StatusCode)
		})
	}
}
