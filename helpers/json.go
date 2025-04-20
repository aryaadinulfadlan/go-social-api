package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aryaadinulfadlan/go-social-api/internal/env"
	"github.com/golang-jwt/jwt/v5"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func ReadFromRequestBody(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing body request")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteToResponseBody(w http.ResponseWriter, statusCode int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(value)
}

func WriteErrorResponse(w http.ResponseWriter, code int, status string, message string) {
	errorResponse := ErrorResponse{
		Code:    code,
		Message: message,
		Status:  status,
	}
	WriteToResponseBody(w, code, errorResponse)
}

func JSONFormatting(data any) {
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonData))
}

func GenerateJWT(user_id string, exp time.Time) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user_id,
		"iat":     time.Now().Unix(),
		"exp":     exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString([]byte(env.Envs.SECRET_KEY))
	if err != nil {
		return "", err
	}
	return token_string, nil
}

func ParseJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(env.Envs.SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}
	return claims, nil
}
