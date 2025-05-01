package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
func TestBasicAuth(t *testing.T) {
	router := SetupTest()
	request := httptest.NewRequest(http.MethodGet, "/v1/basic", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Basic YWRpbnVsOmFkaW51bDEyMw==")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)
}
func TestBasicAuthUnauthorized(t *testing.T) {
	router := SetupTest()
	request := httptest.NewRequest(http.MethodGet, "/v1/basic", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Basic invalid base64")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
