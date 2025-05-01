package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	router := SetupTest()
	request := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	bodyBytes, _ := io.ReadAll(response.Body)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "PONG", string(bodyBytes))
	time.Sleep(time.Second)
}

func TestPingWithRateLimiter(t *testing.T) {
	router := SetupTest()
	for range 5 {
		request := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		response := recorder.Result()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	failedRequest := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	failedRecorder := httptest.NewRecorder()
	router.ServeHTTP(failedRecorder, failedRequest)
	failedResponse := failedRecorder.Result()
	assert.Equal(t, http.StatusTooManyRequests, failedResponse.StatusCode)
	time.Sleep(time.Second)
	successRequest := httptest.NewRequest(http.MethodGet, "/v1/ping", nil)
	successRecorder := httptest.NewRecorder()
	router.ServeHTTP(successRecorder, successRequest)
	successResponse := successRecorder.Result()
	assert.Equal(t, http.StatusOK, successResponse.StatusCode)
}

func TestBasicAuthorized(t *testing.T) {
	router := SetupTest()
	request := httptest.NewRequest(http.MethodGet, "/v1/basic", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic YWRpbnVsOmFkaW51bDEyMw==")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	bodyBytes, _ := io.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Authenticated as Basic Authentication", bodyString)
}

func TestBasicUnauthorized(t *testing.T) {
	router := SetupTest()
	request := httptest.NewRequest(http.MethodGet, "/v1/basic", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic invalid base64")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestBearerAuthorized(t *testing.T) {
	router := SetupTest()
	user_id, _ := uuid.Parse("50b466de-2de4-4e40-bdec-08270f23a8c8")
	token, err := GenerateJWT(user_id.String())
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodGet, "/v1/bearer", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	bodyBytes, _ := io.ReadAll(response.Body)
	bodyString := string(bodyBytes)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Authenticated as Bearer Authentication", bodyString)
}

func TestBearerUnauthorized(t *testing.T) {
	router := SetupTest()
	request := httptest.NewRequest(http.MethodGet, "/v1/bearer", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer invalid token")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
