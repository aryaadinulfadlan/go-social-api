package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/entity/user"
	"github.com/aryaadinulfadlan/go-social-api/internal/shared"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserLogin(t *testing.T) {
	time.Sleep(time.Second) // TO PREVENT RATE LIMITING
	router := SetupTest()
	tests := []struct {
		name               string
		email              string
		password           string
		expectedHttpCode   int
		expectedBodyCode   int
		expectedBodyStatus string
	}{
		{
			name:               "Invalid Email / Password",
			email:              "invalid_email@gmail.com",
			password:           "password123",
			expectedHttpCode:   http.StatusUnauthorized,
			expectedBodyCode:   http.StatusUnauthorized,
			expectedBodyStatus: shared.StatusUnauthorized,
		},
		{
			name:               "In-complete Credentials",
			email:              "princess_diana@gmail.com",
			password:           "",
			expectedHttpCode:   http.StatusBadRequest,
			expectedBodyCode:   http.StatusBadRequest,
			expectedBodyStatus: shared.StatusBadRequest,
		},
		{
			name:               "Valid Credentials",
			email:              "princess_diana@gmail.com",
			password:           "diana123",
			expectedHttpCode:   http.StatusOK,
			expectedBodyCode:   http.StatusOK,
			expectedBodyStatus: shared.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			payload := user.LoginUserPayload{
				Email:    test.email,
				Password: test.password,
			}
			jsonData, _ := json.Marshal(payload)
			requestBody := bytes.NewReader(jsonData)
			request := httptest.NewRequest(http.MethodPost, "/v1/auth/sign-in", requestBody)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedHttpCode, response.StatusCode)
		})
	}
}
func TestUserGetDetail(t *testing.T) {
	time.Sleep(time.Second) // TO PREVENT RATE LIMITING
	router := SetupTest()
	userId, _ := uuid.Parse("e1b4e485-fa48-4d59-8758-e7f988d5cc17")
	tests := []struct {
		name               string
		userId             string
		expectedHttpCode   int
		expectedBodyCode   int
		expectedBodyStatus string
	}{
		{
			name:               "Invalid User ID",
			userId:             "e3488ac6-7012-4d95-a002-663b9a6f879x",
			expectedHttpCode:   http.StatusBadRequest,
			expectedBodyCode:   http.StatusBadRequest,
			expectedBodyStatus: shared.StatusBadRequest,
		},
		{
			name:               "User Doesn't Exist",
			userId:             "e3488ac6-7012-4d95-a002-663b9a6f879a",
			expectedHttpCode:   http.StatusNotFound,
			expectedBodyCode:   http.StatusNotFound,
			expectedBodyStatus: shared.StatusNotFound,
		},
		{
			name:               "Valid User",
			userId:             userId.String(),
			expectedHttpCode:   http.StatusOK,
			expectedBodyCode:   http.StatusOK,
			expectedBodyStatus: shared.StatusOK,
		},
	}
	token, _ := GenerateJWT(userId.String())
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/v1/users/"+test.userId, nil)
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)
			response := recorder.Result()
			assert.Equal(t, test.expectedHttpCode, response.StatusCode)
		})
	}
}
